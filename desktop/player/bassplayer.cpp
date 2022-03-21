#include <QCoreApplication>
#include <QDir>
#include <QFileInfo>
#include <QMap>
#include <QMessageBox>
#include <QUrl>
#include <QtCore>

#include "bassplayer.h"
#include "bass_fx.h"
#include "tags.h"

#if defined(Q_OS_WIN)
#    include "bassasio.h"
#    include "bassmix.h"
#    include "basswasapi.h"
#endif

BassPlayer::BassPlayer()
{
    QString bassPluginsPath = QCoreApplication::applicationDirPath() +
#if defined(Q_OS_MAC)
                              "/../Plugins";
#else
                              "/bassplugins";
#endif
    QDir dir(bassPluginsPath);
    auto plugins = dir.entryList(QStringList() <<
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

BassPlayer::~BassPlayer()
{
#if defined(Q_OS_WIN)
    if (m_asioInitialized)
        BASS_ASIO_Free();
    if (m_wasapiInitialized)
        BASS_WASAPI_Free();
#endif
    BASS_Free();
}

QString BassPlayer::openAudio(const QString &uri)
{
    BASS_ChannelStop(m_hNowPlay);
    BASS_StreamFree(m_hNowPlay);

    if (QFile::exists(uri))
    {
        QString     ext       = QFileInfo(uri).suffix().toLower();
        static const QVector<QString> audioExts = {
            "mp3",
            "mp2",
            "mp1",
            "ogg",
            "wav",
            "aiff",
            "ape",
            "m4a",
            "aac",
            "tta",
            "wma",
            "mp4",
            "m4v",
            "alac",
            "flac",
            "wv",
        };
        if (!audioExts.contains(ext, Qt::CaseInsensitive))
        {
            return "err";
        }
        m_hNowPlay = BASS_StreamCreateFile(false,
#if defined(Q_OS_WIN)
                                           uri.toStdWString().c_str(),
#else
                                           uri.toStdString().c_str(),
#endif
                                           0,
                                           0,
                                           BASS_UNICODE | BASS_SAMPLE_FLOAT | BASS_SAMPLE_FX | BASS_STREAM_DECODE | BASS_STREAM_PRESCAN);
    }
    else if (uri.startsWith("https://", Qt::CaseInsensitive) || uri.startsWith("http://", Qt::CaseInsensitive))
    {
        m_hNowPlay = BASS_StreamCreateURL(
#if defined(Q_OS_WIN)
            uri.toStdWString().c_str(),
#else
            uri.toStdString().c_str(),
#endif
            BASS_UNICODE | BASS_SAMPLE_FLOAT | BASS_SAMPLE_FX | BASS_STREAM_DECODE | BASS_STREAM_PRESCAN,
            0,
            nullptr,
            0);
    }

    if (BASS_ErrorGetCode() != 0)
    {
        return "err";
    }

    if (m_hNowPlay)
    {
        m_hNowPlay = BASS_FX_ReverseCreate(m_hNowPlay, 2, BASS_FX_FREESOURCE | BASS_SAMPLE_FLOAT /*|BASS_SAMPLE_FX*/);
        BASS_ChannelSetAttribute(m_hNowPlay, BASS_ATTRIB_REVERSE_DIR, BASS_FX_RVS_FORWARD);
    }

    m_hReverbFX = BASS_ChannelSetFX(m_hNowPlay, BASS_FX_DX8_REVERB, 1);
    return "ok";
}

void BassPlayer::eqReady()
{
    BASS_BFX_PEAKEQ peakEQ;

    m_hEqFX = BASS_ChannelSetFX(m_hNowPlay, BASS_FX_BFX_PEAKEQ, 2);

    peakEQ.fGain      = 0;
    peakEQ.fQ         = 0;
    peakEQ.fBandwidth = 2.5f;
    peakEQ.lChannel   = BASS_BFX_CHANALL;

    peakEQ.lBand   = 0;
    peakEQ.fCenter = 31;
    BASS_FXSetParameters(m_hEqFX, &peakEQ);

    peakEQ.lBand   = 1;
    peakEQ.fCenter = 62;
    BASS_FXSetParameters(m_hEqFX, &peakEQ);

    peakEQ.lBand   = 2;
    peakEQ.fCenter = 125;
    BASS_FXSetParameters(m_hEqFX, &peakEQ);

    peakEQ.lBand   = 3;
    peakEQ.fCenter = 250;
    BASS_FXSetParameters(m_hEqFX, &peakEQ);

    peakEQ.lBand   = 4;
    peakEQ.fCenter = 500;
    BASS_FXSetParameters(m_hEqFX, &peakEQ);

    peakEQ.lBand   = 5;
    peakEQ.fCenter = 1000;
    BASS_FXSetParameters(m_hEqFX, &peakEQ);

    peakEQ.lBand   = 6;
    peakEQ.fCenter = 2000;
    BASS_FXSetParameters(m_hEqFX, &peakEQ);

    peakEQ.lBand   = 7;
    peakEQ.fCenter = 4000;
    BASS_FXSetParameters(m_hEqFX, &peakEQ);

    peakEQ.lBand   = 8;
    peakEQ.fCenter = 8000;
    BASS_FXSetParameters(m_hEqFX, &peakEQ);

    peakEQ.lBand   = 9;
    peakEQ.fCenter = 16000;
    BASS_FXSetParameters(m_hEqFX, &peakEQ);
}

void BassPlayer::disableEQ()
{
    BASS_ChannelRemoveFX(m_hNowPlay, m_hEqFX);
}

void BassPlayer::setEQ(int id, int gain)
{
    BASS_BFX_PEAKEQ peakEQ;

    peakEQ.lBand = id;
    BASS_FXGetParameters(m_hEqFX, &peakEQ);
    peakEQ.fGain = gain;
    BASS_FXSetParameters(m_hEqFX, &peakEQ);
}

void BassPlayer::setVol(int vol)
{
    float v = (float)vol / 100;
    BASS_ChannelSetAttribute(m_hNowPlay, BASS_ATTRIB_VOL, v);
}

int BassPlayer::getVol()
{
    float vol;
    BASS_ChannelGetAttribute(m_hNowPlay, BASS_ATTRIB_VOL, &vol);
    return (int)(vol * 100);
}

bool BassPlayer::isPlaying()
{
    return (BASS_ChannelIsActive(m_hNowPlay) == BASS_ACTIVE_PLAYING);
}

void BassPlayer::getFFT(float *array)
{
    if (BASS_ChannelIsActive(m_hNowPlay) == BASS_ACTIVE_PLAYING)
        BASS_ChannelGetData(m_hNowPlay, array, BASS_DATA_FFT4096);
}

void BassPlayer::play()
{
    BASS_ChannelPlay(m_hNowPlay, false);
}

void BassPlayer::stop()
{
    BASS_ChannelStop(m_hNowPlay);
    BASS_ChannelSetPosition(m_hNowPlay, 0, BASS_POS_BYTE);
}

void BassPlayer::pause()
{
    BASS_ChannelPause(m_hNowPlay);
}

bool BassPlayer::devInit()
{
    return BASS_Init(-1, 48000, 0, 0, NULL);
}

QString BassPlayer::getTags()
{
    QString tags = TAGS_Read(m_hNowPlay, "%IFV2(%ARTI,%UTF8(%ARTI) - ,)%IFV2(%TITL,%UTF8(%TITL),)");
    if (tags.trimmed().isEmpty())
        return "Show_File_Name";

    return tags;
}

int BassPlayer::getPos()
{
    return (int)(BASS_ChannelGetPosition(m_hNowPlay, BASS_POS_BYTE) * 1000 / BASS_ChannelGetLength(m_hNowPlay, BASS_POS_BYTE));
}

void BassPlayer::setPos(int pos)
{
    BASS_ChannelSetPosition(m_hNowPlay, pos * BASS_ChannelGetLength(m_hNowPlay, BASS_POS_BYTE) / 1000, BASS_POS_BYTE);
}

int BassPlayer::getBitRate()
{
    float time    = BASS_ChannelBytes2Seconds(m_hNowPlay, BASS_ChannelGetLength(m_hNowPlay, BASS_POS_BYTE));
    DWORD len     = BASS_StreamGetFilePosition(m_hNowPlay, BASS_FILEPOS_END);
    int   bitrate = (int)(len / (125 * time) + 0.5);
    return bitrate;
}

int BassPlayer::getFreq()
{
    BASS_CHANNELINFO cInfo;
    BASS_ChannelGetInfo(m_hNowPlay, &cInfo);
    return cInfo.freq;
}

void BassPlayer::setFreq(float freq)
{
    BASS_ChannelSetAttribute(m_hNowPlay, BASS_ATTRIB_FREQ, freq);
}

QString BassPlayer::getNowPlayInfo()
{
    QString           fmt;
    BASS_CHANNELINFO  cInfo;
    BASS_CHANNELINFO *info = &cInfo;
    BASS_ChannelGetInfo(m_hNowPlay, info);

    QMap<DWORD, QString> types = {
        {BASS_CTYPE_STREAM_AIFF, " AIFF"},
        {BASS_CTYPE_STREAM_MP3, " MP3"},
        {BASS_CTYPE_STREAM_MP2, " MP2"},
        {BASS_CTYPE_STREAM_MP1, " MP1"},
        {BASS_CTYPE_STREAM_OGG, " OGG"},
        {BASS_CTYPE_STREAM_WAV_PCM, " Wave PCM"},
        {BASS_CTYPE_STREAM_WAV_FLOAT, " Wave Float Point"},
        //        {BASS_CTYPE_STREAM_APE, " APE"},
        //        {BASS_CTYPE_STREAM_MP4, " MP4"},
        //        {BASS_CTYPE_STREAM_AAC, " AAC"},
        //        {BASS_CTYPE_STREAM_ALAC, " ALAC"},
        //        {BASS_CTYPE_STREAM_TTA, " TTA"},
        //        {BASS_CTYPE_STREAM_FLAC, " FLAC"},
        //        {BASS_CTYPE_STREAM_WMA, " WMA"},
        //        {BASS_CTYPE_STREAM_WMA_MP3, " WMA"},
        //        {BASS_CTYPE_STREAM_WV, " WV"},
    };

    fmt = types.value(info->ctype, "");
    return QString("%1Hz %2Kbps %3%4").arg(info->freq).arg(getBitRate()).arg((info->chans == 1) ? QObject::tr("mono") : QObject::tr("stereo"), fmt);
}

QString BassPlayer::getCurTime()
{
    int totalSec = (int)BASS_ChannelBytes2Seconds(m_hNowPlay, BASS_ChannelGetPosition(m_hNowPlay, BASS_POS_BYTE));
    int minute   = totalSec / 60;
    int second   = totalSec % 60;
    if (second != -1)
    {
        return QString("%1:%2").arg(minute).arg(second, 2, 10, QChar('0'));
    }
    return "0:00";
}

QString BassPlayer::getTotalTime()
{
    int totalSec = (int)BASS_ChannelBytes2Seconds(m_hNowPlay, BASS_ChannelGetLength(m_hNowPlay, BASS_POS_BYTE));
    int minute   = totalSec / 60;
    int second   = totalSec % 60;
    if (second != -1)
    {
        return QString("%1:%2").arg(minute).arg(second, 2, 10, QChar('0'));
    }
    return "0:00";
}

int BassPlayer::getCurTimeMS()
{
    return (int)(BASS_ChannelBytes2Seconds(m_hNowPlay, BASS_ChannelGetPosition(m_hNowPlay, BASS_POS_BYTE)) * 1000);
}

int BassPlayer::getTotalTimeMS()
{
    return (int)(BASS_ChannelBytes2Seconds(m_hNowPlay, BASS_ChannelGetLength(m_hNowPlay, BASS_POS_BYTE)) * 1000);
}

double BassPlayer::getCurTimeSec()
{
    return BASS_ChannelBytes2Seconds(m_hNowPlay, BASS_ChannelGetPosition(m_hNowPlay, BASS_POS_BYTE));
}

double BassPlayer::getTotalTimeSec()
{
    return BASS_ChannelBytes2Seconds(m_hNowPlay, BASS_ChannelGetLength(m_hNowPlay, BASS_POS_BYTE));
}

DWORD BassPlayer::getLevel()
{
    DWORD level = BASS_ChannelGetLevel(m_hNowPlay);
    if (level != (DWORD)-1)
    {
        return level;
    }
    return 0;
}

QString BassPlayer::getFileTotalTime(const QString &fileName)
{
    HSTREAM fileStream = BASS_StreamCreateFile(false, fileName.toStdWString().c_str(), 0, 0, BASS_UNICODE);
    int     totalSec   = (int)BASS_ChannelBytes2Seconds(fileStream, BASS_ChannelGetLength(fileStream, BASS_POS_BYTE));
    BASS_StreamFree(fileStream);
    int minute = totalSec / 60;
    int second = totalSec % 60;
    if (second != -1)
    {
        return QString("%1:%2").arg(minute).arg(second, 2, 10, QChar('0'));
    }
    return "0:00";
}

double BassPlayer::getFileSecond(const QString &fileName)
{
    HSTREAM fileStream = BASS_StreamCreateFile(false, fileName.toStdWString().c_str(), 0, 0, BASS_UNICODE);
    double  totalSec   = BASS_ChannelBytes2Seconds(fileStream, BASS_ChannelGetLength(fileStream, BASS_POS_BYTE));
    BASS_StreamFree(fileStream);
    return totalSec;
}

void BassPlayer::setReverse(bool isEnable)
{
    BASS_ChannelSetAttribute(m_hNowPlay, BASS_ATTRIB_REVERSE_DIR, isEnable ? BASS_FX_RVS_REVERSE : BASS_FX_RVS_FORWARD);
}

void BassPlayer::updateReverb(int value)
{
    BASS_DX8_REVERB p;
    BASS_FXGetParameters(m_hReverbFX, &p);
    p.fReverbMix = 0.012f * (value * value * value);
    BASS_FXSetParameters(m_hReverbFX, &p);
}

BassDriver BassPlayer::getDriver() const
{
    return m_driver;
}

void BassPlayer::setDriver(BassDriver &driver)
{
    m_driver = driver;
}
#if defined(Q_OS_WIN)
bool BassPlayer::asioInit()
{
    if (!BASS_ASIO_Init(-1, 0))
    {
        return false;
    }
    // initialize BASS "no sound" device
    BASS_Init(0, 48000, 0, 0, NULL);
    // create a dummy stream for reserving ASIO channels
    auto dummy = BASS_StreamCreate(2, 48000, BASS_SAMPLE_FLOAT | BASS_STREAM_DECODE, STREAMPROC_DUMMY, NULL);

    // prepare ASIO output channel pairs (up to 4)
    int            a;
    BASS_ASIO_INFO i;
    BASS_ASIO_GetInfo(&i);
    for (a = 0; a < 4; a++)
    {
        BASS_ASIO_CHANNELINFO i, i2;
        if (BASS_ASIO_ChannelGetInfo(FALSE, a * 2, &i) && BASS_ASIO_ChannelGetInfo(FALSE, a * 2 + 1, &i2))
        {
            char name[200];
            sprintf_s(name, "%s + %s", i.name, i2.name);
            // MESS(30 + a, WM_SETTEXT, 0, name);                  // display channel names
            BASS_ASIO_ChannelEnableBASS(FALSE, 0, dummy, TRUE); // enable ASIO channels using the dummy stream
            BASS_ASIO_ChannelPause(FALSE, a * 2);               // not playing anything immediately, so pause the channel
        }
    }

    // start the device using default buffer/latency and 2 threads for parallel processing
    if (!BASS_ASIO_Start(0, 2))
    {
        return false;
    }
    m_asioInitialized = true;
    return true;
}

DWORD BassPlayer::WasapiProc(void *buffer, DWORD length, void *user)
{
    BassPlayer *pThis = (BassPlayer *)user;
    DWORD   c     = BASS_ChannelGetData(pThis->m_mixer, buffer, length);
    if (c == -1)
        c = 0; // an error, no data
    return c;
}

bool BassPlayer::wasapiInit()
{
    // not playing anything via BASS, so don't need an update thread
    BASS_SetConfig(BASS_CONFIG_UPDATEPERIOD, 0);
    // setup BASS - "no sound" device
    BASS_Init(0, 48000, 0, 0, NULL);

    // initialize the default WASAPI device (400ms buffer, 50ms update period, auto-select format)
    if (!BASS_WASAPI_Init(-1, 0, 0, BASS_WASAPI_AUTOFORMAT | BASS_WASAPI_EXCLUSIVE, 0.4, 0.05, &BassPlayer::WasapiProc, (void *)this))
    {
        // exclusive mode failed, try shared mode
        if (!BASS_WASAPI_Init(-1, 0, 0, BASS_WASAPI_AUTOFORMAT, 0.4, 0.05, &BassPlayer::WasapiProc, (void *)this))
        {
            return false;
        }
    }
    BASS_WASAPI_INFO wi;
    BASS_WASAPI_GetInfo(&wi);
    // create a mixer with same format as the output
    m_mixer = BASS_Mixer_StreamCreate(wi.freq, wi.chans, BASS_SAMPLE_FLOAT | BASS_STREAM_DECODE);
    // start the output
    BASS_WASAPI_Start();

    m_wasapiInitialized = true;
    return true;
}
#endif
