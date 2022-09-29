#include <Ice/Ice.h>
#include <typeinfo>
#include <string>
#include <Printer.h>
#include <unistd.h>
 
using namespace std;
using namespace Demo;
 
class PrinterI : public Printer
{
public:
    virtual void Echo(const ::Demo::bytes&, ::Demo::bytes&, const ::Ice::Current& );
};
 
void 
PrinterI::Echo(const ::Demo::bytes& req, ::Demo::bytes& res, const ::Ice::Current&)
{
    string temp(req.begin(), req.end());
    cout << temp.c_str() << endl;
    res = req;
}

int
main(int argc, char* argv[])
{
    int status = 0;
    Ice::CommunicatorPtr ic;
    try
    {
        ic = Ice::initialize(argc, argv);
        Ice::ObjectAdapterPtr adapter =
            ic->createObjectAdapterWithEndpoints("SimplePrinterAdapter", "default -p 10000");
        Ice::ObjectPtr object = new PrinterI;
        adapter->add(object, ic->stringToIdentity("SimplePrinter"));
        adapter->activate();
        ic->waitForShutdown();
    }
    catch(const std::exception& e)
    {
        cerr << e.what() << endl;
        return 1;
    }
    if (ic) {
        try {
            ic->destroy();
        } catch (const Ice::Exception& e) {
            cerr << e << endl;
            status = 1;
        }
    }
    return status;
}
