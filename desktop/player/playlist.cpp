#include <QStandardPaths>

#include "playlist.h"
#include "ui_playlist.h"

PlayList::PlayList(Player *player, QWidget *parent) : QWidget(parent), ui(new Ui::PlayList), m_player(player)
{
    ui->setupUi(this);
    ui->playListTable->horizontalHeader()->setStretchLastSection(true);
    setWindowFlags(Qt::WindowTitleHint | Qt::CustomizeWindowHint);
    setFixedSize(width(), height());

    connect(ui->searchEdit, SIGNAL(returnPressed()), this, SLOT(on_searchButton_clicked()));

    readFromFile(QStandardPaths::writableLocation(QStandardPaths::AppDataLocation) + "/default.hpl");
}

PlayList::~PlayList()
{
    auto dirPath = QStandardPaths::writableLocation(QStandardPaths::AppDataLocation);
    QDir dir(dirPath);
    if (!dir.exists())
    {
        dir.mkdir(dirPath);
    }
    saveToFile(dirPath + "/default.hpl");
    delete ui;
}

bool PlayList::fixSuffix(const QString &fileName)
{
    QString     ext       = QFileInfo(fileName).suffix().toLower();
    QStringList audioExts = {
        "mp3",
        "mp1",
        "ogg",
        "ape",
        "m4a",
        "aac",
        "tta",
        "wma",
        "mp2",
        "wav",
        "aiff",
        "mp4",
        "m4v",
        "alac",
        "flac",
        "wv",
    };
    return audioExts.contains(ext);
}

bool PlayList::isEmpty()
{
    return m_trackList.isEmpty();
}

void PlayList::add(const QString &fileName)
{
    if (fixSuffix(fileName))
    {
        if ((int)(m_player->getFileSecond(fileName)) >= m_lengthFilter)
        {
            m_trackList.append(fileName);
            m_timeList.append(m_player->getFileTotalTime(fileName));
        }
    }
    tableUpdate();
}

void PlayList::insert(int index, const QString &fileName)
{
    if (fixSuffix(fileName))
    {
        if ((int)(m_player->getFileSecond(fileName)) >= m_lengthFilter)
        {
            if (index < m_curIndex)
                ++m_curIndex;
            m_trackList.insert(index, fileName);
            m_timeList.insert(index, m_player->getFileTotalTime(fileName));
        }
    }
    tableUpdate();
}

void PlayList::remove(int index)
{
    if (index <= m_curIndex && index > -1)
        --m_curIndex;
    m_trackList.removeAt(index);
    m_timeList.removeAt(index);
    tableUpdate();
}

void PlayList::clearAll()
{
    m_trackList.clear();
    m_timeList.clear();
    m_curIndex = 0;
    tableUpdate();
}

int PlayList::getLength()
{
    return m_trackList.length();
}

int PlayList::getIndex()
{
    if (!m_trackList.isEmpty())
    {
        return m_curIndex;
    }
    return -1;
}

QString PlayList::next(bool isLoop)
{
    if (!m_trackList.isEmpty())
    {
        if (isLoop)
        {
            if (m_curIndex < m_trackList.length() - 1)
            {
                ++m_curIndex;
            }
            else
            {
                m_curIndex = 0;
            }
            ui->playListTable->selectRow(m_curIndex);
            tableUpdate();
            return m_trackList[m_curIndex];
        }
        if (m_curIndex < m_trackList.length() - 1)
        {
            ++m_curIndex;
            ui->playListTable->selectRow(m_curIndex);
            tableUpdate();
            return m_trackList[m_curIndex];
        }
        return "stop";
    }
    return "";
}

QString PlayList::previous(bool isLoop)
{
    if (!m_trackList.isEmpty())
    {
        if (isLoop)
        {
            if (m_curIndex == 0)
            {
                m_curIndex = m_trackList.length() - 1;
            }
            else
            {
                --m_curIndex;
            }
            ui->playListTable->selectRow(m_curIndex);
            tableUpdate();
            return m_trackList[m_curIndex];
        }
        if (m_curIndex > 0)
        {
            --m_curIndex;
            ui->playListTable->selectRow(m_curIndex);
            tableUpdate();
            return m_trackList[m_curIndex];
        }
        return "stop";
    }
    return "";
}

const QString &PlayList::playIndex(int index)
{
    m_curIndex = index;
    ui->playListTable->selectRow(m_curIndex);
    tableUpdate();
    return m_trackList[m_curIndex];
}

const QString &PlayList::getFileNameForIndex(int index)
{
    return m_trackList[index];
}

const QString &PlayList::getCurFile()
{
    return m_trackList[m_curIndex];
}

QString PlayList::playLast()
{
    if (!m_trackList.isEmpty())
    {
        m_curIndex = m_trackList.length() - 1;
        ui->playListTable->selectRow(m_curIndex);
        tableUpdate();
        return m_trackList[m_curIndex];
    }
    return "";
}

void PlayList::tableUpdate()
{
    ui->playListTable->clear();
    ui->playListTable->setRowCount(getLength());
    int count = m_trackList.size();
    for (int i = 0; i < count; i++)
    {
        QString   fileName = m_trackList[i];
        QFileInfo fileInfo(fileName);

        QTableWidgetItem *item     = new QTableWidgetItem(fileInfo.fileName());
        QTableWidgetItem *timeItem = new QTableWidgetItem(m_timeList[i]);

        if (i == m_curIndex)
        {
            item->setBackground(QColor(128, 255, 0, 128));
            timeItem->setBackground(QColor(128, 255, 0, 128));
        }

        ui->playListTable->setItem(i, 0, item);
        ui->playListTable->setItem(i, 1, timeItem);
    }
}

void PlayList::on_deleteButton_clicked()
{
    QItemSelectionModel *selectionModel = ui->playListTable->selectionModel();
    QModelIndexList      selected       = selectionModel->selectedIndexes();
    QMap<int, int>       rowMap;
    for (const auto &index : selected)
    {
        rowMap.insert(index.row(), 0);
    }

    QMapIterator<int, int> rowMapIterator(rowMap);
    rowMapIterator.toBack();
    while (rowMapIterator.hasPrevious())
    {
        rowMapIterator.previous();
        remove(rowMapIterator.key());
    }
}

void PlayList::dragEnterEvent(QDragEnterEvent *event)
{
    event->acceptProposedAction();
}

void PlayList::dropEvent(QDropEvent *event)
{
    auto urls = event->mimeData()->urls();
    if (urls.isEmpty())
        return;

    for (const auto &u : urls)
    {
        auto fileName = u.toLocalFile();
        if (fixSuffix(fileName))
        {
            add(fileName);
        }
    }
}

void PlayList::on_playListTable_cellDoubleClicked(int row, int)
{
    m_curIndex = row;
    emit callPlayer();
    tableUpdate();
}

void PlayList::on_clearButton_clicked()
{
    clearAll();
}

void PlayList::on_insertButton_clicked()
{
    int index = ui->playListTable->currentRow();
    if (index < 0)
        index = 0;
    QStringList fileNames = QFileDialog::getOpenFileNames(
        this,
        tr("Insert before selected item"),
        0,
        tr("Audio file (*.mp3 *.mp2 *.mp1 *.wav *.aiff *.ogg *.ape *.mp4 *.m4a *.m4v *.aac *.alac *.tta *.flac *.wma *.wv)"));
    int count = fileNames.size();
    for (int i = 0; i < count; i++)
    {
        QString fileName = fileNames[i];
        insert(index + i, fileName);
    }
}

void PlayList::on_addButton_clicked()
{
    QStringList fileNames = QFileDialog::getOpenFileNames(
        this,
        tr("Add audio"),
        0,
        tr("Audio file (*.mp3 *.mp2 *.mp1 *.wav *.aiff *.ogg *.ape *.mp4 *.m4a *.m4v *.aac *.alac *.tta *.flac *.wma *.wv)"));
    for (const auto &fileName : fileNames)
    {
        add(fileName);
    }
}

void PlayList::on_searchButton_clicked()
{
    if (!m_trackList.isEmpty() && !ui->searchEdit->text().isEmpty())
    {
        int resultIndex = -1;
        int count       = m_trackList.size();
        for (int i = 0; i < count; i++)
        {
            QString   fileName = m_trackList[i];
            QFileInfo fileInfo(fileName);

            if (ui->isCaseSensitive->isChecked())
            {
                if (fileInfo.fileName().indexOf(ui->searchEdit->text()) > -1)
                {
                    resultIndex = i;
                    break;
                }
            }
            else
            {
                if (fileInfo.fileName().toLower().indexOf(ui->searchEdit->text().toLower()) > -1)
                {
                    resultIndex = i;
                    break;
                }
            }
        }

        if (resultIndex != -1)
        {
            ui->playListTable->selectRow(resultIndex);
        }
        else
        {
            QMessageBox::information(this, tr("Sorry"), tr("Cannot find."));
        }
    }
    else if (ui->searchEdit->text().isEmpty())
    {
        QMessageBox::question(this, tr("Hello"), tr("What are you looking for?"));
    }
    else
    {
        QMessageBox::question(this, tr("What's this"), tr("Why I should search in the empty list?"));
    }
}

void PlayList::on_searchNextButton_clicked()
{
    if (!m_trackList.isEmpty() && !ui->searchEdit->text().isEmpty())
    {
        int resultIndex = -1;
        int start       = ui->playListTable->currentRow() + 1;
        int count       = m_trackList.size();

        if (start < count)
            for (int i = start; i < count; i++)
            {
                QString   fileName = m_trackList[i];
                QFileInfo fileInfo(fileName);

                if (ui->isCaseSensitive->isChecked())
                {
                    if (fileInfo.fileName().indexOf(ui->searchEdit->text()) > -1)
                    {
                        resultIndex = i;
                        break;
                    }
                }
                else
                {
                    if (fileInfo.fileName().toLower().indexOf(ui->searchEdit->text().toLower()) > -1)
                    {
                        resultIndex = i;
                        break;
                    }
                }
            }

        if (resultIndex != -1)
        {
            ui->playListTable->selectRow(resultIndex);
        }
        else
        {
            QMessageBox::information(0, tr("Searching done!"), tr("All things are searched."));
        }
    }
    else if (ui->searchEdit->text().isEmpty())
    {
        QMessageBox::question(0, tr("Hello"), tr("What are you looking for?"));
    }
    else
    {
        QMessageBox::question(this, tr("What's this"), tr("Why I should search in the empty list?"));
    }
}

void PlayList::on_setLenFilButton_clicked()
{
    bool ok;
    int  set = QInputDialog::getInt(0,
                                   tr("Minimum Playback Length"),
                                   tr("Audio files smaller than this length will not be accepted \n unit: seconds"),
                                   m_lengthFilter,
                                   0,
                                   2147483647,
                                   1,
                                   &ok);
    if (ok)
        m_lengthFilter = set;
}

void PlayList::saveToFile(const QString &fileName)
{
    QFile file(fileName);
    if (file.open(QIODevice::WriteOnly | QIODevice::Truncate))
    {
        QDataStream stream(&file);
        stream << (quint32)0x61727487 << m_trackList << m_timeList << m_curIndex;
        file.close();
    }
}

void PlayList::readFromFile(const QString &fileName)
{
    QFile file(fileName);
    if (file.open(QIODevice::ReadOnly))
    {
        QDataStream stream(&file);
        quint32     magic;
        stream >> magic;
        if (magic == 0x61727487)
        {
            stream >> m_trackList;
            stream >> m_timeList;
            stream >> m_curIndex;
        }
        file.close();
        tableUpdate();
    }
}
