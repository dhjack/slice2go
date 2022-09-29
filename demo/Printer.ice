module Demo
{
	sequence<byte> bytes;
	interface Printer
	{
		 void Echo(bytes req, out bytes res);
	};
};
