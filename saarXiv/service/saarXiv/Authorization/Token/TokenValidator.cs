using System;
using System.Linq;
using System.Security.Cryptography;
using System.Security.Cryptography.X509Certificates;
using System.Threading.Tasks;
using saarXiv.Models;

namespace saarXiv.Authorization.Token
{
    public class TokenValidator
    {

     public async Task<TokenResult> ValidateAsync(Paper paper, string downloadToken)
     {

         try
         {
             byte[] tokenBytes = Convert.FromBase64String(downloadToken);
             var cert = new X509Certificate2(tokenBytes);

             if (!cert.SignatureAlgorithm.Value.Equals("1.2.840.10045.4.3.2"))
             {
                 return TokenResult.Failed("Invalid signature algorithm");
             }

             if (!cert.GetPublicKey().SequenceEqual(paper.Author.PublicKey))
             {
                 return TokenResult.Failed($"Token not valid for key {Convert.ToBase64String(paper.Author.PublicKey)}");
             }

             if (!cert.GetNameInfo(X509NameType.SimpleName, false).Equals(paper.SaarXivID))
             {
                 return TokenResult.Failed($"Token not valid for {paper.SaarXivID}");
             }

             var chain = new X509Chain
             {
                 ChainPolicy =
                 {
                     RevocationMode = X509RevocationMode.NoCheck,
                     VerificationFlags = X509VerificationFlags.AllowUnknownCertificateAuthority
                 }
             };
             if (!chain.Build(cert))
             {
                 return TokenResult.Failed("Invalid Signature");
             }
         }
         catch (FormatException)
         {
             return TokenResult.Failed("Invalid Format");
         }
         catch (CryptographicException)
         {
             return TokenResult.Failed("Invalid Token");
         }
         catch (NullReferenceException)
         {
             return TokenResult.Failed("Invalid Token");
         }

         return TokenResult.Success;
     }
    }
}