using System.Collections.Concurrent;
using System.Text.Json;
using Application.EventBus;
using Application.EventBus.Events;
using Application.Factories.Game;
using Application.Services.Game;
using Microsoft.AspNetCore.SignalR;
using Microsoft.Extensions.Logging;
using SIPSorcery.Net;

namespace Infrastructure.Hubs;

public class ArenaWebRtcService
{
    private readonly ConcurrentDictionary<string, RTCPeerConnection> _peerConnections = new();
    private readonly ConcurrentDictionary<string, RTCDataChannel> _dataChannels = new();
    private readonly ConcurrentDictionary<string, List<RTCIceCandidateInit>> _iceBuffer = new();
    
    private readonly IArenaService _arenaService;
    private readonly IArenaMessageHandlerFactory _handlerFactory;
    private readonly ILogger<ArenaWebRtcService> _logger;

    public ArenaWebRtcService(
        IArenaService arenaService,
        IArenaMessageHandlerFactory handlerFactory,
        ILogger<ArenaWebRtcService> logger,
        IEventBus eventBus)
    {
        _arenaService = arenaService;
        _handlerFactory = handlerFactory;
        _logger = logger;

        // Подписка на событие SnapshotReady через EventBus
        eventBus.Subscribe<SnapshotReadyEvent>(async e =>
        {
            if (_dataChannels.TryGetValue(e.ClientId, out var dc) &&
                dc.readyState == RTCDataChannelState.open)
            {
                var options = new JsonSerializerOptions
                {
                    PropertyNamingPolicy = JsonNamingPolicy.CamelCase
                };
                
                var packet = new { type = "snapshotResponse", payload = e.State };
                var json = JsonSerializer.Serialize(packet, options);
                dc.send(json);
            }
        });
    }

    // -------------------------------------------
    // Handle WebRTC Offer from client
    // -------------------------------------------
    public async Task<string> HandleOfferFromClientAsync(
        string clientId,
        string sdpOffer,
        IClientProxy clientProxy)
    {
        _logger.LogInformation($"[{nameof(ArenaWebRtcService)}]: Offer received from {clientId}");

        var pcConfig = new RTCConfiguration
        {
            iceServers = new[] { new RTCIceServer { urls = "stun:stun.l.google.com:19302" } }.ToList()
        };

        var pc = new RTCPeerConnection(pcConfig);
        _peerConnections[clientId] = pc; // Регистрируем сразу (до await)

        // --- Подписки на события PeerConnection ---
        pc.onicecandidate += candidate =>
        {
            if (candidate is not null)
            {
                var candidateJson = JsonSerializer.Serialize(candidate);
                clientProxy.SendAsync("ServerIceCandidate", candidateJson);
            }
        };

        pc.ondatachannel += dc =>
        {
            _logger.LogInformation($"[{nameof(ArenaWebRtcService)}]: DataChannel '{dc.label}' created for {clientId}");
            _dataChannels[clientId] = dc;

            dc.onopen += async () =>
            {
                _logger.LogInformation($"[{nameof(ArenaWebRtcService)}]: DataChannel '{dc.label}' opened for {clientId}");
                await SendSnapshotToClientAsync(clientId);
            };

            dc.onmessage += async (_, _, data) =>
            {
                var json = System.Text.Encoding.UTF8.GetString(data);
                await HandleClientMessageAsync(clientId, json);
            };

            dc.onclose += () =>
            {
                _logger.LogInformation($"[{nameof(ArenaWebRtcService)}]: DataChannel '{dc.label}' closed for {clientId}");
                _dataChannels.TryRemove(clientId, out _);
            };
        };

        // --- SDP Setup ---
        var remoteDesc = new RTCSessionDescriptionInit
        {
            type = RTCSdpType.offer,
            sdp = sdpOffer
        };

        pc.setRemoteDescription(remoteDesc);

        var answer = pc.createAnswer();
        await pc.setLocalDescription(answer);

        // --- Применяем отложенные ICE кандидаты ---
        if (_iceBuffer.TryRemove(clientId, out var buffered))
        {
            foreach (var c in buffered)
                pc.addIceCandidate(c);
        }

        _logger.LogInformation($"[{nameof(ArenaWebRtcService)}]: Answer created for {clientId}");
        return answer.sdp;
    }

    // -------------------------------------------
    // ICE Candidate from client
    // -------------------------------------------
    public void AddIceCandidate(string clientId, string candidateJson)
    {
        var candidate = JsonSerializer.Deserialize<RTCIceCandidateInit>(candidateJson);
        if (candidate is null) return;

        if (_peerConnections.TryGetValue(clientId, out var pc))
        {
            pc.addIceCandidate(candidate);
        }
        else
        {
            var buffer = _iceBuffer.GetOrAdd(clientId, _ => new List<RTCIceCandidateInit>());
            buffer.Add(candidate);
            _logger.LogDebug($"[{nameof(ArenaWebRtcService)}]: Buffered ICE candidate for {clientId}");
        }
    }

    // -------------------------------------------
    // Graceful Close
    // -------------------------------------------
    public void ClosePeer(string clientId)
    {
        if (_dataChannels.TryRemove(clientId, out var dc))
        {
            try { dc.close(); }
            catch
            {
                // ignored
            }
        }

        if (_peerConnections.TryRemove(clientId, out var pc))
        {
            try { pc.close(); }
            catch
            {
                // ignored
            }
        }

        _logger.LogInformation($"[{nameof(ArenaWebRtcService)}]: Closed peer {clientId}");
    }

    // -------------------------------------------
    // Handle messages from client
    // -------------------------------------------
    private async Task HandleClientMessageAsync(string clientId, string message)
    {
        var doc = JsonDocument.Parse(message);

        if (!doc.RootElement.TryGetProperty("type", out var t)) return;
        var type = t.GetString() ?? string.Empty;
        var payload = doc.RootElement.TryGetProperty("payload", out var p) ? p : default;

        var handler = _handlerFactory.CreateHandler(type);
        if (handler is null)
        {
            _logger.LogWarning($"[{nameof(ArenaWebRtcService)}]: No handler for message type '{type}'");
            return;
        }

        await handler.HandleAsync(clientId, payload);
    }

    // -------------------------------------------
    // Send snapshot immediately via data channel
    // -------------------------------------------
    private async Task SendSnapshotToClientAsync(string clientId)
    {
        if (!_dataChannels.TryGetValue(clientId, out var dc)) return;
        if (dc.readyState != RTCDataChannelState.open) return;

        var state = await _arenaService.GetStateAsync(clientId);
        var packet = new { type = "snapshotResponse", payload = state };
        var json = JsonSerializer.Serialize(packet);
        dc.send(json);

        _logger.LogDebug($"[{nameof(ArenaWebRtcService)}]: Snapshot sent to {clientId}");
    }
}