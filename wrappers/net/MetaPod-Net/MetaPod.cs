using System;
using System.IO;
using System.Runtime.InteropServices;
using System.Text;
using System.Threading.Tasks;

namespace MetaPod_Net
{
   public static class MetaPod
       {
           /// <summary>
           /// This method attempts to create a new executable containing the provided payload.
           /// The payload is injected into the provided input file/template.
           /// If successful the new MetaPod portable executable is written to the output path.
           /// </summary>
           /// <param name="inputFile">The path to the digitally signed portable executable that metadata will be written to.</param>
           /// <param name="payload">The metadata that will be written to the portable executable.</param>
           /// <param name="outputFile">The output file that will be generated upon successful completion.</param>
           /// <exception cref="FileNotFoundException"></exception>
           /// <exception cref="ArgumentNullException"></exception>
           /// <exception cref="MetaPodException"></exception>
           public static async Task Create(string inputFile, string payload, string outputFile)
           {
               if (!File.Exists(inputFile))
               {
                   throw new FileNotFoundException($"Unable to locate {inputFile}.");
               }
               if (string.IsNullOrWhiteSpace(payload))
               {
                   throw new ArgumentNullException($"Payload string cannot be null or empty.");
               }
               var inputBytes = await File.ReadAllBytesAsync(inputFile);
               var payloadBytes = Encoding.UTF8.GetBytes(payload);
               var output = IntPtr.Zero;
               try
               {
                   var errorCode = 0;
                   var outputSize  = NativeWrapper.Create(inputBytes, inputBytes.Length, payloadBytes, ref output, ref errorCode);
                   if (errorCode > 0)
                   {
                       throw new MetaPodException(GetErrorCodeMessage(errorCode));
                   }
                   var outputBuffer = new byte[outputSize];
                   Marshal.Copy(output, outputBuffer, 0, outputSize);
                   using (var fs = new FileStream(outputFile, FileMode.Create, FileAccess.Write, FileShare.None, outputBuffer.Length, true))
                   {
                       await fs.WriteAsync(outputBuffer, 0, outputBuffer.Length);
                   }
               }
               finally
               {
                   Marshal.FreeHGlobal(output);
               }
           }
   
           /// <summary>
           /// This method attempts to create a new executable containing the provided payload.
           /// The payload is injected into the provided input file/template.
           /// </summary>
           /// <param name="inputFile">The raw bytes of the digitally signed portable executable that metadata will be written to.</param>
           /// <param name="payload">The metadata that will be written to the portable executable.</param>
           /// <returns>The new MetaPod portable executable.</returns>
           /// <exception cref="ArgumentNullException"></exception>
           /// <exception cref="MetaPodException"></exception>
           public static Span<byte> Create(Span<byte> inputFile, string payload)
           {
               if (inputFile.Length == 0)
               {
                   throw new ArgumentNullException($"Input file bytes cannot be zero.");
               }
               if (string.IsNullOrWhiteSpace(payload))
               {
                   throw new ArgumentNullException($"Payload string cannot be null or empty.");
               }
               var payloadBytes = Encoding.UTF8.GetBytes(payload);
               var output = IntPtr.Zero;
               try
               {
                   var outputSize = 0;
                   
                   var errorCode  = NativeWrapper.Create(inputFile.ToArray(), inputFile.Length, payloadBytes, ref output, ref outputSize);
                   if (errorCode > 0)
                   {
                       throw new MetaPodException(GetErrorCodeMessage(errorCode));
                   }
                   var outputBuffer = new byte[outputSize];
                   Marshal.Copy(output, outputBuffer, 0, outputSize);
                   return new Span<byte>(outputBuffer);
               }
               finally
               {
                   Marshal.FreeHGlobal(output);
               }
           }
   
           /// <summary>
           /// This method attempts to open and read a MetaPod portable executable from a provided input file path.
           /// If no payload is found a  <see cref="MetaPodException"/> will be thrown.
           /// </summary>
           /// <param name="inputFile">The path to your previous created MetaPod portable executable.</param>
           /// <returns></returns>
           /// <exception cref="FileNotFoundException"></exception>
           /// <exception cref="MetaPodException"></exception>
           public static async Task<string> Open(string inputFile)
           {
               if (!File.Exists(inputFile))
               {
                   throw new FileNotFoundException($"Unable to locate {inputFile}.");
               }
               var inputBytes = await File.ReadAllBytesAsync(inputFile);
               var output = IntPtr.Zero;
               var payloadSize = 0;
               try
               {
                   var errorCode = NativeWrapper.Open(inputBytes, inputBytes.Length, ref output, ref payloadSize);
                   if (errorCode > 0)
                   {
                       throw new MetaPodException(GetErrorCodeMessage(errorCode));
                   }
                   var byteArray = new byte[payloadSize];
                   Marshal.Copy(output, byteArray, 0, payloadSize);
                   return Encoding.UTF8.GetString(byteArray);
               }
               finally
               {
                   Marshal.FreeHGlobal(output);
               }
           }
   
           /// <summary>
           /// This method attempts to open and read a MetaPod portable executable from a provided <see cref="Span{T}"/>
           /// If no payload is found a  <see cref="MetaPodException"/> will be thrown.
           /// </summary>
           /// <param name="inputFile"></param>
           /// <returns>The MetaPod payload as a string.</returns>
           public static string Open(Span<byte> inputFile)
           {
               var output = IntPtr.Zero;
               var payloadSize = 0;
               try
               {
                   var errorCode = NativeWrapper.Open(inputFile.ToArray(), inputFile.Length, ref output, ref payloadSize);
                   if (errorCode > 0)
                   {
                       throw new MetaPodException(GetErrorCodeMessage(errorCode));
                   }
                   var byteArray = new byte[payloadSize];
                   Marshal.Copy(output, byteArray, 0, payloadSize);
                   return Encoding.UTF8.GetString(byteArray);
               }
               finally
               {
                   Marshal.FreeHGlobal(output);
               }
           }
   
           /// <summary>
           /// Uses a provided error code to safely fetch the relevant error message.
           /// </summary>
           /// <param name="errorCode"></param>
           /// <returns>The MetaPod error message.</returns>
           private static string GetErrorCodeMessage(int errorCode)
           {
               var output = IntPtr.Zero;
               try
               {
                   var messageSize = NativeWrapper.GetErrorCodeMessage(errorCode, ref output);
                   var byteArray = new byte[messageSize];
                   Marshal.Copy(output, byteArray, 0, messageSize);
                   return Encoding.UTF8.GetString(byteArray);
               }
               finally
               {
                   Marshal.FreeHGlobal(output);
               }
           }
       }
}