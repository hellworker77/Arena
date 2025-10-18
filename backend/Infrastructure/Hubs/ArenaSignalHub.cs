using Application.Services.Game;
using Microsoft.AspNetCore.SignalR;
using Microsoft.Extensions.Logging;

namespace Infrastructure.Hubs;

public class ArenaSignalHub(ArenaWebRtcService webRtc,
    IArenaService arenaService,
    ILogger<ArenaSignalHub> logger): Hub
{
    public async Task OfferToServer(string sdpOffer)
    {
        var callerId = Context.ConnectionId;
        
        var sdpAnswer = await webRtc.HandleOfferFromClientAsync(callerId, sdpOffer, Clients.Caller);
        
        await Clients.Caller.SendAsync("AnswerFromServer", sdpAnswer);
    }

    public Task SendIceCandidateToServer(string candidateJson)
    {
        var callerId = Context.ConnectionId;
        webRtc.AddIceCandidate(callerId, candidateJson);
        return Task.CompletedTask;
    }

    public override async Task OnConnectedAsync()
    {
        var connectionId = Context.ConnectionId;
        logger.LogInformation($"[{nameof(ArenaSignalHub)}]: Connected: {Context.ConnectionId}");

        await arenaService.InitializeSessionAsync(connectionId);
        
        await base.OnConnectedAsync();
    }
    
    public override async Task OnDisconnectedAsync(Exception? exception)
    {
        var connectionId = Context.ConnectionId;
        logger.LogInformation($"[{nameof(ArenaSignalHub)}]: Disconnected: {Context.ConnectionId}");
        webRtc.ClosePeer(Context.ConnectionId);

        await arenaService.CloseSessionAsync(connectionId);
        
        await base.OnDisconnectedAsync(exception);
    }
}