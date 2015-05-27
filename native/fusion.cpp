// This code based on the OpenCV tutorial
#include "fusion.h"
#include <opencv2/opencv.hpp>
#include <fstream>

using namespace cv;
using namespace std;

void runFusion(const char* output, const char* input)
{
    try
    {
        ifstream lst;
        lst.open(input);
        if (!lst.is_open()) {
            return;
        }
        vector<Mat> imageset;
        string filename;
        for (lst >> filename; lst; lst >> filename)
        {
            imageset.push_back(imread(filename));
        }
        lst.close();
        Mat fusion;
        Ptr<MergeMertens> mertens = createMergeMertens();
        mertens->process(imageset, fusion);
        imwrite(output, fusion * 255);
    } catch (...) {

    }
}