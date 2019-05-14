using System;

namespace MetaPod_Net
{
    [Serializable]
    public class MetaPodException : Exception
    {
        public MetaPodException()
        {

        }

        public MetaPodException(string message)
            : base($"MetaPod Error: {message}")
        {

        }
    }
}