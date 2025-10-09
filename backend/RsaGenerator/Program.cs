using System.Security.Cryptography;
using System.IO;

using var rsaLocal = RSA.Create(2048);

var privatePem = ExportPrivateKeyPem(rsaLocal);
var publicPem = ExportPublicKeyPem(rsaLocal);

File.WriteAllText("private.pem", privatePem);
File.WriteAllText("public.pem", publicPem);

Console.WriteLine("Keys generated and saved to private.pem and public.pem");

string ExportPrivateKeyPem(RSA rsa)
{
    var key = rsa.ExportRSAPrivateKey();
    
    return "-----BEGIN PRIVATE KEY-----\n" +
           Convert.ToBase64String(key, Base64FormattingOptions.InsertLineBreaks) +
           "\n-----END PRIVATE KEY-----";
}

string ExportPublicKeyPem(RSA rsa)
{
    var key = rsa.ExportRSAPublicKey();
    
    return "-----BEGIN PUBLIC KEY-----\n" +
           Convert.ToBase64String(key, Base64FormattingOptions.InsertLineBreaks) +
           "\n-----END PUBLIC KEY-----";
}