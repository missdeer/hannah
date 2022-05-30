#ifndef PLAYER_H
#define PLAYER_H
#include <QString>

#include "bass.h"

enum BassDriver
{
    BD_Default,
    BD_ASIO,
    BD_WASAPI,
};

class BassPlayer
{
public:
    Q_DISABLE_COPY_MOVE(BassPlayer)
    BassPlayer();
    ~BassPlayer();
    bool    devInit();
    [[nodiscard]] QString getTags();
    void                  setVol(int vol);
    [[nodiscard]] int     getVol();
    QString openAudio(const QString &uri);
    void    play();
    void    pause();
    void    stop();
    [[nodiscard]] int     getPos();
    void                  setPos(int pos);
    [[nodiscard]] int     getBitRate();
    void                  getFFT(float *array);
    [[nodiscard]] bool    isPlaying();
    [[nodiscard]] int     getFreq();
    void    setFreq(float freq);
    void    eqReady();
    void    disableEQ();
    void                  setEQ(int id, int gain);
    [[nodiscard]] QString getNowPlayInfo();
    [[nodiscard]] QString getTotalTime();
    [[nodiscard]] QString getCurTime();
    [[nodiscard]] int     getCurTimeMS() const;
    [[nodiscard]] int     getTotalTimeMS() const;
    [[nodiscard]] double  getCurTimeSec() const;
    [[nodiscard]] double  getTotalTimeSec() const;
    [[nodiscard]] DWORD   getLevel() const;
    [[nodiscard]] QString getFileTotalTime(const QString &fileName);
    [[nodiscard]] double  getFileSecond(const QString &fileName);
    void    setReverse(bool isEnable) const;
    void    updateReverb(int value) const;
    void    setJumpPoint(double timeFrom, double timeTo);
    void    removeJumpPoint();

    [[nodiscard]] BassDriver getDriver() const;
    void       setDriver(BassDriver &driver);
    // WASAPI function
    static DWORD CALLBACK WasapiProc(void *buffer, DWORD length, void *user);

private:
    HSTREAM m_hNowPlay {0};
    HFX     m_hEqFX {0};
    HFX     m_hReverbFX {0};
    HSYNC   m_hJumping {0};
    bool    m_bPlayNextEnable {true};
#if defined(Q_OS_WIN)
    HSTREAM m_mixer {0};
    bool    m_asioInitialized {false};
    bool    m_wasapiInitialized {false};
    bool    asioInit();
    bool    wasapiInit();
#endif
    BassDriver m_driver {BD_Default};
};

inline BassPlayer *gBassPlayer = nullptr;

#endif // PLAYER_H
