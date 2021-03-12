#include <QCoreApplication>
#include <QDir>
#include <QFileInfo>
#include <QMessageBox>

#include "player.h"

#include "bass_fx.h"
#include "tags.h"

Player::Player()
{
    //加载解码插件
    QString bassPluginsPath = QCoreApplication::applicationDirPath() + "/bassplugins";
    QDir    dir(bassPluginsPath);
    auto    plugins = dir.entryList(QStringList() <<
#if defined(Q_OS_WIN)
                                 "*.dll"
#elif defined(Q_OS_MAC)
                                 "*.dylib"
#else
                                 "*.so"
#endif
    );
    for (auto &p : plugins)
    {
        QString plugin = QDir::toNativeSeparators(bassPluginsPath + "/" + p);
        BASS_PluginLoad(plugin.toStdString().c_str(), 0);
    }
}

Player::~Player()
{
    BASS_Free();
}

QString Player::openFile(QString fileName)
{
    //打开新文件前停止播放
    QString ext = QFileInfo(fileName).suffix();
    if (ext.compare("mp3", Qt::CaseInsensitive) == 0 || ext.compare("mp2", Qt::CaseInsensitive) == 0 ||
        ext.compare("mp1", Qt::CaseInsensitive) == 0 || ext.compare("wav", Qt::CaseInsensitive) == 0 ||
        ext.compare("ogg", Qt::CaseInsensitive) == 0 || ext.compare("aiff", Qt::CaseInsensitive) == 0 ||
        ext.compare("ape", Qt::CaseInsensitive) == 0 || ext.compare("mp4", Qt::CaseInsensitive) == 0 ||
        ext.compare("m4a", Qt::CaseInsensitive) == 0 || ext.compare("m4v", Qt::CaseInsensitive) == 0 ||
        ext.compare("aac", Qt::CaseInsensitive) == 0 || ext.compare("alac", Qt::CaseInsensitive) == 0 ||
        ext.compare("tta", Qt::CaseInsensitive) == 0 || ext.compare("flac", Qt::CaseInsensitive) == 0 ||
        ext.compare("wma", Qt::CaseInsensitive) == 0 || ext.compare("wv", Qt::CaseInsensitive) == 0)
    {
        BASS_ChannelStop(nowPlay);
        BASS_StreamFree(nowPlay);
        nowPlay = BASS_StreamCreateFile(false,
                                        fileName.toStdWString().c_str(),
                                        0,
                                        0,
                                        BASS_UNICODE | BASS_SAMPLE_FLOAT | BASS_SAMPLE_FX | BASS_STREAM_DECODE | BASS_STREAM_PRESCAN);

        if (BASS_ErrorGetCode() != 0)
        {
            return "err"; //打开文件失败
        }

        if (nowPlay) //如果流创建成功
        {
            //创建反向流
            nowPlay = BASS_FX_ReverseCreate(nowPlay, 2, BASS_FX_FREESOURCE | BASS_SAMPLE_FLOAT /*|BASS_SAMPLE_FX*/);
            //转为正向播放
            BASS_ChannelSetAttribute(nowPlay, BASS_ATTRIB_REVERSE_DIR, BASS_FX_RVS_FORWARD);
        }

        //创建混响效果
        reverbFX = BASS_ChannelSetFX(nowPlay, BASS_FX_DX8_REVERB, 1);
        //此时默认开启混响……需要关闭
        return "ok";
    }
    else
    {
        return "err";
    }
}

void Player::eqReady()
{
    BASS_BFX_PEAKEQ peakEQ;

    eqFX = BASS_ChannelSetFX(nowPlay, BASS_FX_BFX_PEAKEQ, 2);

    peakEQ.fGain      = 0;
    peakEQ.fQ         = 0;
    peakEQ.fBandwidth = 2.5f;
    peakEQ.lChannel   = BASS_BFX_CHANALL;

    peakEQ.lBand   = 0;
    peakEQ.fCenter = 31;
    BASS_FXSetParameters(eqFX, &peakEQ);

    peakEQ.lBand   = 1;
    peakEQ.fCenter = 62;
    BASS_FXSetParameters(eqFX, &peakEQ);

    peakEQ.lBand   = 2;
    peakEQ.fCenter = 125;
    BASS_FXSetParameters(eqFX, &peakEQ);

    peakEQ.lBand   = 3;
    peakEQ.fCenter = 250;
    BASS_FXSetParameters(eqFX, &peakEQ);

    peakEQ.lBand   = 4;
    peakEQ.fCenter = 500;
    BASS_FXSetParameters(eqFX, &peakEQ);

    peakEQ.lBand   = 5;
    peakEQ.fCenter = 1000;
    BASS_FXSetParameters(eqFX, &peakEQ);

    peakEQ.lBand   = 6;
    peakEQ.fCenter = 2000;
    BASS_FXSetParameters(eqFX, &peakEQ);

    peakEQ.lBand   = 7;
    peakEQ.fCenter = 4000;
    BASS_FXSetParameters(eqFX, &peakEQ);

    peakEQ.lBand   = 8;
    peakEQ.fCenter = 8000;
    BASS_FXSetParameters(eqFX, &peakEQ);

    peakEQ.lBand   = 9;
    peakEQ.fCenter = 16000;
    BASS_FXSetParameters(eqFX, &peakEQ);
}

void Player::disableEQ()
{
    BASS_ChannelRemoveFX(nowPlay, eqFX);
}

void Player::setEQ(int id, int gain)
{
    BASS_BFX_PEAKEQ peakEQ; //均衡器参数结构
    // id对应均衡器段编号
    peakEQ.lBand = id;
    BASS_FXGetParameters(eqFX, &peakEQ);
    peakEQ.fGain = gain;                 //-15~15，EQ参数
    BASS_FXSetParameters(eqFX, &peakEQ); //变更参数
}

void Player::setVol(int vol)
{
    //设置音量，0最小，100最大
    float v = (float)vol / 100;
    BASS_ChannelSetAttribute(nowPlay, BASS_ATTRIB_VOL, v);
}

int Player::getVol()
{
    //取得音量，0最小，100最大
    float vol;
    BASS_ChannelGetAttribute(nowPlay, BASS_ATTRIB_VOL, &vol);
    return (int)(vol * 100);
}

bool Player::isPlaying()
{
    if (BASS_ChannelIsActive(nowPlay) == BASS_ACTIVE_PLAYING)
        return true;
    else
        return false;
}

void Player::getFFT(float *array)
{
    if (BASS_ChannelIsActive(nowPlay) == BASS_ACTIVE_PLAYING)
        BASS_ChannelGetData(nowPlay, array, BASS_DATA_FFT4096);
}

void Player::play()
{
    BASS_ChannelPlay(nowPlay, false);
}

void Player::stop()
{
    BASS_ChannelStop(nowPlay);
    BASS_ChannelSetPosition(nowPlay, 0, BASS_POS_BYTE);
}

void Player::pause()
{
    BASS_ChannelPause(nowPlay);
}

//初始化音频设备
bool Player::devInit()
{
    return BASS_Init(-1, 48000, 0, 0, NULL);
}

QString Player::getTags()
{
    // QString tags = TAGS_Read(nowPlay, "%IFV2(%ARTI,%UTF8(%ARTI),未知艺术家) - %IFV2(%TITL,%UTF8(%TITL),无标题)");

    //有一些音频把艺术家写到了标题里
    //很少见到只有艺术家没有标题的音频
    //故修改为下列表达式，若只有艺术家没有标题会是“艺术家 - ”的形式……喵
    //（末尾为" - "应该删去3个字符？）
    QString tags = TAGS_Read(nowPlay, "%IFV2(%ARTI,%UTF8(%ARTI) - ,)%IFV2(%TITL,%UTF8(%TITL),)");
    if (tags.trimmed().isEmpty())
        return "Show_File_Name"; //如果标签是空字符，直接显示文件名
    else
        return tags; //返回标签
}

int Player::getPos()
{
    //返回当前播放位置，取值范围0~1000
    return (int)(BASS_ChannelGetPosition(nowPlay, BASS_POS_BYTE) * 1000 / BASS_ChannelGetLength(nowPlay, BASS_POS_BYTE));
}

void Player::setPos(int pos)
{
    //跳转进度到指定位置，0~1000
    BASS_ChannelSetPosition(nowPlay, pos * BASS_ChannelGetLength(nowPlay, BASS_POS_BYTE) / 1000, BASS_POS_BYTE);
}

//取得音频比特率
int Player::getBitRate()
{
    float time    = BASS_ChannelBytes2Seconds(nowPlay, BASS_ChannelGetLength(nowPlay, BASS_POS_BYTE)); // 播放时间
    DWORD len     = BASS_StreamGetFilePosition(nowPlay, BASS_FILEPOS_END);                             // 文件长度
    int   bitrate = (int)(len / (125 * time) + 0.5);                                                   // 比特率/编码率 (Kbps)
    return bitrate;
}

//取得音频采样率
int Player::getFreq()
{
    BASS_CHANNELINFO cInfo;
    BASS_ChannelGetInfo(nowPlay, &cInfo);
    return cInfo.freq;
}

//设置音频采样率
void Player::setFreq(float freq)
{
    BASS_ChannelSetAttribute(nowPlay, BASS_ATTRIB_FREQ, freq);
}

QString Player::getNowPlayInfo()
{
    QString           fmt;
    BASS_CHANNELINFO  cInfo;
    BASS_CHANNELINFO *info = &cInfo;
    BASS_ChannelGetInfo(nowPlay, info);

    if (info->ctype == BASS_CTYPE_STREAM_AIFF)
        fmt = " AIFF";
    else if (info->ctype == BASS_CTYPE_STREAM_MP3)
        fmt = " MP3";
    else if (info->ctype == BASS_CTYPE_STREAM_MP2)
        fmt = " MP2";
    else if (info->ctype == BASS_CTYPE_STREAM_MP1)
        fmt = " MP1";
    else if (info->ctype == BASS_CTYPE_STREAM_OGG)
        fmt = " OGG";
    else if (info->ctype == BASS_CTYPE_STREAM_WAV_PCM)
        fmt = " Wave PCM";
    else if (info->ctype == BASS_CTYPE_STREAM_WAV_FLOAT)
        fmt = QString::fromUtf8(" Wave Float Point");
    //    else if (info->ctype == BASS_CTYPE_STREAM_APE)
    //        fmt = QString::fromUtf8(" APE");
    //    else if (info->ctype == BASS_CTYPE_STREAM_MP4)
    //        fmt = QString::fromUtf8(" MP4");
    //    else if (info->ctype == BASS_CTYPE_STREAM_AAC)
    //        fmt = QString::fromUtf8(" AAC");
    //    else if (info->ctype == BASS_CTYPE_STREAM_ALAC)
    //        fmt = QString::fromUtf8(" ALAC");
    //    else if (info->ctype == BASS_CTYPE_STREAM_TTA)
    //        fmt = QString::fromUtf8(" TTA");
    //    else if (info->ctype == BASS_CTYPE_STREAM_FLAC)
    //        fmt = QString::fromUtf8(" FLAC");
    //    else if (info->ctype == BASS_CTYPE_STREAM_WMA)
    //        fmt = QString::fromUtf8(" WMA");
    //    else if (info->ctype == BASS_CTYPE_STREAM_WMA_MP3)
    //        fmt = QString::fromUtf8(" WMA");
    //    else if (info->ctype == BASS_CTYPE_STREAM_WV)
    //        fmt = QString::fromUtf8(" WV");

    if (info->chans == 1)
        return QString::number(info->freq) + "Hz " + QString::number(this->getBitRate()) + "Kbps " + QObject::tr("mono") + fmt;
    else
        return QString::number(info->freq) + "Hz " + QString::number(this->getBitRate()) + "Kbps " + QObject::tr("stereo") + fmt;
}

QString Player::getCurTime()
{
    int totalSec = (int)BASS_ChannelBytes2Seconds(nowPlay, BASS_ChannelGetPosition(nowPlay, BASS_POS_BYTE));
    int minute   = totalSec / 60;
    int second   = totalSec % 60;
    if (second != -1)
    {
        if (second < 10)
            return QString::number(minute) + ":0" + QString::number(second);
        else
            return QString::number(minute) + ":" + QString::number(second);
    }
    else
    {
        return "0:00";
    }
}

QString Player::getTotalTime()
{
    int totalSec = (int)BASS_ChannelBytes2Seconds(nowPlay, BASS_ChannelGetLength(nowPlay, BASS_POS_BYTE));
    int minute   = totalSec / 60;
    int second   = totalSec % 60;
    if (second != -1)
    {
        if (second < 10)
            return QString::number(minute) + ":0" + QString::number(second);
        else
            return QString::number(minute) + ":" + QString::number(second);
    }
    else
    {
        return "0:00";
    }
}

int Player::getCurTimeMS()
{
    return (int)(BASS_ChannelBytes2Seconds(nowPlay, BASS_ChannelGetPosition(nowPlay, BASS_POS_BYTE)) * 1000);
}

int Player::getTotalTimeMS()
{
    return (int)(BASS_ChannelBytes2Seconds(nowPlay, BASS_ChannelGetLength(nowPlay, BASS_POS_BYTE)) * 1000);
}

double Player::getCurTimeSec()
{
    return BASS_ChannelBytes2Seconds(nowPlay, BASS_ChannelGetPosition(nowPlay, BASS_POS_BYTE));
}

double Player::getTotalTimeSec()
{
    return BASS_ChannelBytes2Seconds(nowPlay, BASS_ChannelGetLength(nowPlay, BASS_POS_BYTE));
}

DWORD Player::getLevel()
{
    //返回声道电平，低16位左声道，高16位右声道
    DWORD level = BASS_ChannelGetLevel(nowPlay);
    if (level != -1)
    {
        return level;
    }
    else
    {
        return 0;
    }
}

QString Player::getFileTotalTime(QString fileName)
{
    HSTREAM fileStream = BASS_StreamCreateFile(false, fileName.toStdWString().c_str(), 0, 0, BASS_UNICODE);
    int     totalSec   = (int)BASS_ChannelBytes2Seconds(fileStream, BASS_ChannelGetLength(fileStream, BASS_POS_BYTE));
    BASS_StreamFree(fileStream);
    int minute = totalSec / 60;
    int second = totalSec % 60;
    if (second != -1)
    {
        if (second < 10)
            return QString::number(minute) + ":0" + QString::number(second);
        else
            return QString::number(minute) + ":" + QString::number(second);
    }
    else
    {
        return "0:00";
    }
}

double Player::getFileSecond(QString fileName)
{
    HSTREAM fileStream = BASS_StreamCreateFile(false, fileName.toStdWString().c_str(), 0, 0, BASS_UNICODE);
    double  totalSec   = BASS_ChannelBytes2Seconds(fileStream, BASS_ChannelGetLength(fileStream, BASS_POS_BYTE));
    BASS_StreamFree(fileStream);
    return totalSec;
}

//更改播放方向（false正、true反）
void Player::setReverse(bool isEnable)
{
    if (isEnable)
        BASS_ChannelSetAttribute(nowPlay, BASS_ATTRIB_REVERSE_DIR, BASS_FX_RVS_REVERSE);
    else
        BASS_ChannelSetAttribute(nowPlay, BASS_ATTRIB_REVERSE_DIR, BASS_FX_RVS_FORWARD);
}

//更新混响效果，参数取值范围：-20~0
void Player::updateReverb(int value)
{
    BASS_DX8_REVERB p;
    BASS_FXGetParameters(reverbFX, &p);
    p.fReverbMix = 0.012f * (value * value * value); //参数取值范围：-96~0
    BASS_FXSetParameters(reverbFX, &p);
}
