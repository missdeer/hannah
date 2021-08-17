
#ifndef _ShadowPower_ID3V2PIC___
#define _ShadowPower_ID3V2PIC___
#define _CRT_SECURE_NO_WARNINGS
#ifndef NULL
#    define NULL 0
#endif
#include <cstdio>
#include <cstdlib>
#include <cstring>
#include <memory.h>

using namespace std;

namespace spID3
{
    struct ID3V2Header
    {
        char          identi[3];
        unsigned char major;
        unsigned char revsion;
        unsigned char flags;
        unsigned char size[4];
    };

    struct ID3V2FrameHeader
    {
        char          FrameId[4];
        unsigned char size[4];
        unsigned char flags[2];
    };

    struct ID3V22FrameHeader
    {
        char          FrameId[3];
        unsigned char size[3];
    };

    inline unsigned char *pPicData     = 0;
    inline int            picLength    = 0;
    inline char           picFormat[4] = {};

    inline bool verificationPictureFormat(char *data)
    {
        // supported format: JPEG/PNG/BMP/GIF
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

    //安全释放内存
    inline void freePictureData()
    {
        if (pPicData)
        {
            delete pPicData;
        }
        pPicData  = 0;
        picLength = 0;
        memset(&picFormat, 0, 4);
    }

    inline bool loadPictureData(const char *inFilePath)
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

        ID3V2Header id3v2h;

        if (strncmp(id3v2h.identi, "ID3", 3) != 0)
        {
            fclose(fp);
            fp = NULL;
            return false;
        }

        int tagTotalLength =
            (id3v2h.size[0] & 0x7f) * 0x200000 + (id3v2h.size[1] & 0x7f) * 0x4000 + (id3v2h.size[2] & 0x7f) * 0x80 + (id3v2h.size[3] & 0x7f);

        if (id3v2h.major == 3 || id3v2h.major == 4)
        {
            ID3V2FrameHeader id3v2fh;
            memset(&id3v2fh, 0, 10);

            bool hasExtendedHeader = ((id3v2h.flags >> 6 & 0x1) == 1);

            if (hasExtendedHeader)
            {
                unsigned char extendedHeaderSize[4] = {};
                memset(&extendedHeaderSize, 0, 4);
                fread(&extendedHeaderSize, 4, 1, fp);

                int extendedHeaderLength =
                    extendedHeaderSize[0] * 0x1000000 + extendedHeaderSize[1] * 0x10000 + extendedHeaderSize[2] * 0x100 + extendedHeaderSize[3];

                fseek(fp, extendedHeaderLength, SEEK_CUR);
            }

            fread(&id3v2fh, 10, 1, fp);
            int curDataLength = 10;
            while ((strncmp(id3v2fh.FrameId, "APIC", 4) != 0))
            {
                if (curDataLength > tagTotalLength)
                {
                    fclose(fp);
                    fp = NULL;
                    return false;
                }

                int frameLength = id3v2fh.size[0] * 0x1000000 + id3v2fh.size[1] * 0x10000 + id3v2fh.size[2] * 0x100 + id3v2fh.size[3];
                fseek(fp, frameLength, SEEK_CUR);
                memset(&id3v2fh, 0, 10);
                fread(&id3v2fh, 10, 1, fp);
                curDataLength += frameLength + 10;
            }

            int frameLength = id3v2fh.size[0] * 0x1000000 + id3v2fh.size[1] * 0x10000 + id3v2fh.size[2] * 0x100 + id3v2fh.size[3];

            int nonPicDataLength = 0;
            fseek(fp, 1, SEEK_CUR);
            nonPicDataLength++;

            char tempData[20]   = {};
            char mimeType[20]   = {};
            int  mimeTypeLength = 0;

            fread(&tempData, 20, 1, fp);
            fseek(fp, -20, SEEK_CUR);

            strcpy(mimeType, tempData);
            mimeTypeLength = strlen(mimeType) + 1;
            fseek(fp, mimeTypeLength, SEEK_CUR);
            nonPicDataLength += mimeTypeLength;

            fseek(fp, 1, SEEK_CUR);
            nonPicDataLength++;

            int temp = 0;
            fread(&temp, 1, 1, fp);
            nonPicDataLength++;
            while (temp)
            {
                fread(&temp, 1, 1, fp);
                nonPicDataLength++;
            }

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
            picLength = frameLength - nonPicDataLength;
            pPicData  = new unsigned char[picLength];
            memset(pPicData, 0, picLength);
            fread(pPicData, picLength, 1, fp);
            //------------------------
            fclose(fp);
        }
        else if (id3v2h.major == 2)
        {
            // ID3v2.2
            ID3V22FrameHeader id3v2fh;
            memset(&id3v2fh, 0, 6);
            fread(&id3v2fh, 6, 1, fp);
            int curDataLength = 6;
            while ((strncmp(id3v2fh.FrameId, "PIC", 3) != 0))
            {
                if (curDataLength > tagTotalLength)
                {
                    fclose(fp);
                    fp = NULL;
                    return false;
                }
                int frameLength = id3v2fh.size[0] * 0x10000 + id3v2fh.size[1] * 0x100 + id3v2fh.size[2];
                fseek(fp, frameLength, SEEK_CUR);
                memset(&id3v2fh, 0, 6);
                fread(&id3v2fh, 6, 1, fp);
                curDataLength += frameLength + 6;
            }

            int frameLength = id3v2fh.size[0] * 0x10000 + id3v2fh.size[1] * 0x100 + id3v2fh.size[2];

            int nonPicDataLength = 0;
            fseek(fp, 1, SEEK_CUR);
            nonPicDataLength++;

            char imageType[4] = {};
            memset(&imageType, 0, 4);
            fread(&imageType, 3, 1, fp);
            nonPicDataLength += 3;

            fseek(fp, 1, SEEK_CUR);
            nonPicDataLength++;

            int temp = 0;
            fread(&temp, 1, 1, fp);
            nonPicDataLength++;
            while (temp)
            {
                fread(&temp, 1, 1, fp);
                nonPicDataLength++;
            }
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
            picLength = frameLength - nonPicDataLength;
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
            return false;
        }
        return true;
    }

    inline int getPictureLength()
    {
        return picLength;
    }

    inline unsigned char *getPictureDataPtr()
    {
        return pPicData;
    }

    inline char *getPictureFormat()
    {
        return picFormat;
    }

    inline bool writePictureDataToFile(const char *outFilePath)
    {
        FILE *fp = NULL;
        if (picLength > 0)
        {
            fp = fopen(outFilePath, "wb");
            if (fp)
            {
                fwrite(pPicData, picLength, 1, fp);
                fclose(fp);
                return true;
            }
        }
        return false;
    }

    inline bool extractPicture(const char *inFilePath, const char *outFilePath)
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
} // namespace spID3
#endif
