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
#include <memory.h>
#include <cstring>

typedef unsigned char byte;
using namespace std;

namespace spFLAC {
	//Flac元数据块头部结构体定义
	struct FlacMetadataBlockHeader
	{
		byte flag;		//标志位，高1位：是否为最后一个数据块，低7位：数据块类型
		byte length[3];	//数据块长度，不含数据块头部
	};

	byte *pPicData = 0;		//指向图片数据的指针
	int picLength = 0;		//存放图片数据长度
	char picFormat[4] = {};	//存放图片数据的格式（扩展名）

	//检测图片格式，参数1：数据，返回值：是否成功（不是图片则失败）
	bool verificationPictureFormat(char *data)
	{
		//支持格式：JPEG/PNG/BMP/GIF
		byte jpeg[2] = { 0xff, 0xd8 };
		byte png[8] = { 0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a };
		byte gif[6] = { 0x47, 0x49, 0x46, 0x38, 0x39, 0x61 };
		byte gif2[6] = { 0x47, 0x49, 0x46, 0x38, 0x37, 0x61 };
		byte bmp[2] = { 0x42, 0x4d };
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
	void freePictureData()
	{
		if (pPicData)
		{
			delete pPicData;
		}
		pPicData = 0;
		picLength = 0;
		memset(&picFormat, 0, 4);
	}

	//将图片提取到内存，参数1：文件路径，成功返回true
	bool loadPictureData(const char *inFilePath)
	{
		freePictureData();
		FILE *fp = NULL;
		fp = fopen(inFilePath, "rb");
		if (!fp)						//如果打开失败
		{
			fp = NULL;
			return false;
		}
		fseek(fp, 0, SEEK_SET);			//设文件流指针到文件头部
		byte magic[4] = {};				//存放校验数据
		memset(&magic, 0, 4);
		fread(&magic, 4, 1, fp);			//读入校验数据
		byte fLaC[4] = { 0x66, 0x4c, 0x61, 0x43 };
		if (memcmp(&magic, &fLaC, 4) == 0)
		{
			//数据校验正确，文件类型为Flac
			FlacMetadataBlockHeader fmbh;	//创建Flac元数据块头部结构体
			memset(&fmbh, 0, 4);			//清空内存
			fread(&fmbh, 4, 1, fp);			//读入头部数据
			//计算数据块长度，不含头部
			int blockLength = fmbh.length[0] * 0x10000 + fmbh.length[1] * 0x100 + fmbh.length[2];
			int loopCount = 0;	//循环计数，防死
			while ((fmbh.flag & 0x7f) != 6)
			{
				//如果数据类型不是图片，此处循环执行
				loopCount++;
				if (loopCount > 40)
				{
					//循环40次没有遇到末尾就直接停止
					fclose(fp);
					fp = NULL;
					return false;					//可能文件不正常
				}
				fseek(fp, blockLength, SEEK_CUR);	//跳过数据块
				if ((fmbh.flag & 0x80) == 0x80)
				{
					//已经是最后一个数据块了，仍然不是图片
					fclose(fp);
					fp = NULL;
					return false;					//没有找到图片数据
				}
				//取得下一数据块头部
				memset(&fmbh, 0, 4);				//清空内存
				fread(&fmbh, 4, 1, fp);				//读入头部数据
				blockLength = fmbh.length[0] * 0x10000 + fmbh.length[1] * 0x100 + fmbh.length[2];//计算数据块长度
			}
			//此时已到图片数据块

			int nonPicDataLength = 0;				//非图片数据长度
			fseek(fp, 4, SEEK_CUR);					//信仰之跃
			nonPicDataLength += 4;
			char nextJumpLength[4];					//下次要跳的长度
			fread(&nextJumpLength, 4, 1, fp);		//读取安全跳跃距离
			nonPicDataLength += 4;
			int jumpLength = nextJumpLength[0] * 0x1000000 + nextJumpLength[1] * 0x10000 + nextJumpLength[2] * 0x100 + nextJumpLength[3];//计算数据块长度
			fseek(fp, jumpLength, SEEK_CUR);		//Let's Jump!!
			nonPicDataLength += jumpLength;
			fread(&nextJumpLength, 4, 1, fp);
			nonPicDataLength += 4;
			jumpLength = nextJumpLength[0] * 0x1000000 + nextJumpLength[1] * 0x10000 + nextJumpLength[2] * 0x100 + nextJumpLength[3];
			fseek(fp, jumpLength, SEEK_CUR);		//Let's Jump too!!
			nonPicDataLength += jumpLength;
			fseek(fp, 20, SEEK_CUR);				//信仰之跃
			nonPicDataLength += 20;

			//非主流情况检测+获得文件格式
			char tempData[20] = {};
			memset(tempData, 0, 20);
			fread(&tempData, 8, 1, fp);
			fseek(fp, -8, SEEK_CUR);	//回到原位
			//判断40次，一位一位跳到文件头
			bool ok = false;			//是否正确识别出文件头
			for (int i = 0; i < 40; i++)
			{
				//校验文件头
				if (verificationPictureFormat(tempData))
				{
					ok = true;
					break;
				}
				else
				{
					//如果校验失败尝试继续向后校验
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
				return false;			//无法识别的数据
			}

			//-----抵达图片数据区-----
			picLength = blockLength - nonPicDataLength;		//计算图片数据长度
			pPicData = new byte[picLength];					//动态分配图片数据内存空间
			memset(pPicData, 0, picLength);					//清空图片数据内存
			fread(pPicData, picLength, 1, fp);				//得到图片数据
			//------------------------
			fclose(fp);										//操作已完成，关闭文件。
		}
		else
		{
			//校验失败，不是Flac
			fclose(fp);
			fp = NULL;
			freePictureData();
			return false;
		}
		return true;
	}

	//取得图片数据的长度
	int getPictureLength()
	{
		return picLength;
	}

	//取得指向图片数据的指针
	byte *getPictureDataPtr()
	{
		return pPicData;
	}

	//取得图片数据的扩展名（指针）
	char *getPictureFormat()
	{
		return picFormat;
	}

	bool writePictureDataToFile(const char *outFilePath)
	{
		FILE *fp = NULL;
		if (picLength > 0)
		{
			fp = fopen(outFilePath, "wb");		//打开目标文件
			if (fp)								//打开成功
			{
				fwrite(pPicData, picLength, 1, fp);	//写入文件
				fclose(fp);							//关闭
				return true;
			}
			else
			{
				return false;						//文件打开失败
			}
		}
		else
		{
			return false;						//没有图像数据
		}
	}

	//提取图片文件，参数1：输入文件，参数2：输出文件，返回值：是否成功
	bool extractPicture(const char *inFilePath, const char *outFilePath)
	{
		if (loadPictureData(inFilePath))	//如果取得图片数据成功
		{
			if (writePictureDataToFile(outFilePath))
			{
				return true;				//文件写出成功
			}
			else
			{
				return false;				//文件写出失败
			}
		}
		else
		{
			return false;					//无图片数据
		}
		freePictureData();
	}
}
#endif
