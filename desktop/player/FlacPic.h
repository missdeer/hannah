/*
FLAC标签图片提取库 Ver 1.0
从FLAC文件中稳定、快捷、高效、便捷地提取出图片数据
支持BMP、JPEG、PNG、GIF图片格式
可将图片数据提取到文件或内存中，并能安全地释放内存
使用方式与ID3v2版本相同
ShadowPower 于2014/8/1 夜间
*/

#ifndef _ShadowPower_FLACPIC___
#define _ShadowPower_FLACPIC___
#define _CRT_SECURE_NO_WARNINGS
#ifndef NULL
#define NULL 0
#endif
#include <cstdio>
#include <cstdlib>
#include <cstring>
#include <memory.h>

using namespace std;

namespace spFLAC
{
    struct FlacMetadataBlockHeader
    {
        unsigned char flag;
        unsigned char length[3];
    };

    unsigned char *pPicData     = 0;
    int            picLength    = 0;
    char           picFormat[4] = {};

    bool verificationPictureFormat(char *data)
    {
        unsigned char jpeg[2] = {0xff, 0xd8};
        unsigned char png[8]  = {0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a};
        unsigned char gif[6]  = {0x47, 0x49, 0x46, 0x38, 0x39, 0x61};
        unsigned char gif2[6] = {0x47, 0x49, 0x46, 0x38, 0x37, 0x61};
        unsigned char bmp[2]  = {0x42, 0x4d};
        memset(&picFormat, 0, 4);
        if (memcmp(data, &jpeg, 2) == 0)
        {
            strcpy(picFormat, "jpg");
        }
        else if (memcmp(data, &png, 8) == 0)
        {
            strcpy(picFormat, "png");
        }
        else if (memcmp(data, &gif, 6) == 0 || memcmp(data, &gif2, 6) == 0)
        {
            strcpy(picFormat, "gif");
        }
        else if (memcmp(data, &bmp, 2) == 0)
        {
            strcpy(picFormat, "bmp");
        }
        else
        {
            return false;
        }

        return true;
    }

    void freePictureData()
    {
        if (pPicData)
        {
            delete pPicData;
        }
        pPicData  = 0;
        picLength = 0;
        memset(&picFormat, 0, 4);
    }

    bool loadPictureData(const char *inFilePath)
    {
        freePictureData();
        FILE *fp = NULL;
        fp       = fopen(inFilePath, "rb");
        if (!fp)
        {
            fp = NULL;
            return false;
        }
        fseek(fp, 0, SEEK_SET);
        unsigned char magic[4] = {};
        memset(&magic, 0, 4);
        fread(&magic, 4, 1, fp);
        unsigned char fLaC[4] = {0x66, 0x4c, 0x61, 0x43};
        if (memcmp(&magic, &fLaC, 4) == 0)
        {
            FlacMetadataBlockHeader fmbh;
            memset(&fmbh, 0, 4);
            fread(&fmbh, 4, 1, fp);
            int blockLength = fmbh.length[0] * 0x10000 + fmbh.length[1] * 0x100 + fmbh.length[2];
            int loopCount   = 0;
            while ((fmbh.flag & 0x7f) != 6)
            {
                loopCount++;
                if (loopCount > 40)
                {
                    fclose(fp);
                    fp = NULL;
                    return false;
                }
                fseek(fp, blockLength, SEEK_CUR);
                if ((fmbh.flag & 0x80) == 0x80)
                {
                    fclose(fp);
                    fp = NULL;
                    return false;
                }
                memset(&fmbh, 0, 4);
                fread(&fmbh, 4, 1, fp);
                blockLength = fmbh.length[0] * 0x10000 + fmbh.length[1] * 0x100 + fmbh.length[2]; //计算数据块长度
            }

            int nonPicDataLength = 0;
            fseek(fp, 4, SEEK_CUR);
            nonPicDataLength += 4;
            char nextJumpLength[4];
            fread(&nextJumpLength, 4, 1, fp);
            nonPicDataLength += 4;
            int jumpLength =
                nextJumpLength[0] * 0x1000000 + nextJumpLength[1] * 0x10000 + nextJumpLength[2] * 0x100 + nextJumpLength[3]; //计算数据块长度
            fseek(fp, jumpLength, SEEK_CUR);                                                                                 // Let's Jump!!
            nonPicDataLength += jumpLength;
            fread(&nextJumpLength, 4, 1, fp);
            nonPicDataLength += 4;
            jumpLength = nextJumpLength[0] * 0x1000000 + nextJumpLength[1] * 0x10000 + nextJumpLength[2] * 0x100 + nextJumpLength[3];
            fseek(fp, jumpLength, SEEK_CUR); // Let's Jump too!!
            nonPicDataLength += jumpLength;
            fseek(fp, 20, SEEK_CUR);
            nonPicDataLength += 20;

            char tempData[20] = {};
            memset(tempData, 0, 20);
            fread(&tempData, 8, 1, fp);
            fseek(fp, -8, SEEK_CUR);
            bool ok = false;
            for (int i = 0; i < 40; i++)
            {
                if (verificationPictureFormat(tempData))
                {
                    ok = true;
                    break;
                }
                else
                {
                    fseek(fp, 1, SEEK_CUR);
                    nonPicDataLength++;
                    fread(&tempData, 8, 1, fp);
                    fseek(fp, -8, SEEK_CUR);
                }
            }

            if (!ok)
            {
                fclose(fp);
                fp = NULL;
                freePictureData();
                return false;
            }

            picLength = blockLength - nonPicDataLength;
            pPicData  = new unsigned char[picLength];
            memset(pPicData, 0, picLength);
            fread(pPicData, picLength, 1, fp);
            //------------------------
            fclose(fp);
        }
        else
        {
            fclose(fp);
            fp = NULL;
            freePictureData();
            return false;
        }
        return true;
    }

    int getPictureLength()
    {
        return picLength;
    }

    unsigned char *getPictureDataPtr()
    {
        return pPicData;
    }

    char *getPictureFormat()
    {
        return picFormat;
    }

    bool writePictureDataToFile(const char *outFilePath)
    {
        if (picLength > 0)
        {
            FILE *fp = fopen(outFilePath, "wb");
            if (fp)
            {
                fwrite(pPicData, picLength, 1, fp);
                fclose(fp);
                return true;
            }
        }
        return false;
    }

    bool extractPicture(const char *inFilePath, const char *outFilePath)
    {
        if (loadPictureData(inFilePath))
        {
            if (writePictureDataToFile(outFilePath))
            {
                freePictureData();
                return true;
            }
            freePictureData();
        }
        return false;
    }
} // namespace spFLAC
#endif
