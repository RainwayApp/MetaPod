using System;
using System.Runtime.InteropServices;

namespace MetaPod_Net
{
    internal static class NativeWrapper
    {
        private const string LibraryName = "metapod64.so";

        [DllImport(LibraryName, CharSet = CharSet.Unicode, CallingConvention = CallingConvention.Cdecl)]
        [return: MarshalAs(UnmanagedType.LPStr)]
         internal static extern string Open(byte[] portableExecutalbe, int length, ref int errorCode);

         [DllImport(LibraryName, CharSet = CharSet.Unicode, CallingConvention = CallingConvention.Cdecl)]
        [return: MarshalAs(UnmanagedType.LPStr)]
        internal static extern string GetErrorCodeMessage(int errorCode);
        
        
        [DllImport(LibraryName, CharSet = CharSet.Unicode, CallingConvention = CallingConvention.Cdecl)]
        internal static extern int Create(byte[] template, int length, byte[] payload, ref IntPtr output, ref int errorCode);
    }
}