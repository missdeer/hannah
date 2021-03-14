#include "playlist.h"
#include "ui_playlist.h"

PlayList::PlayList(Player *player, QWidget *parent) : QWidget(parent), ui(new Ui::PlayList), m_player(player)
{
    ui->setupUi(this);
    ui->playListTable->horizontalHeader()->setStretchLastSection(true);
    setWindowFlags(Qt::WindowTitleHint | Qt::CustomizeWindowHint);
    setFixedSize(width(), height());

    connect(ui->searchEdit, SIGNAL(returnPressed()), this, SLOT(on_searchButton_clicked()));

    readFromFile(QCoreApplication::applicationDirPath() + "/PlayList.sdpl");
}

PlayList::~PlayList()
{
    saveToFile(QCoreApplication::applicationDirPath() + "/PlayList.sdpl");
    delete ui;
}

bool PlayList::fixSuffix(const QString &fileName)
{
    QFileInfo fi(fileName);
    QString   ext = fi.suffix();
    return (ext.compare("mp3", Qt::CaseInsensitive) == 0 || ext.compare("mp2", Qt::CaseInsensitive) == 0 ||
            ext.compare("mp1", Qt::CaseInsensitive) == 0 || ext.compare("wav", Qt::CaseInsensitive) == 0 ||
            ext.compare("ogg", Qt::CaseInsensitive) == 0 || ext.compare("aiff", Qt::CaseInsensitive) == 0 ||
            ext.compare("ape", Qt::CaseInsensitive) == 0 || ext.compare("mp4", Qt::CaseInsensitive) == 0 ||
            ext.compare("m4a", Qt::CaseInsensitive) == 0 || ext.compare("m4v", Qt::CaseInsensitive) == 0 ||
            ext.compare("aac", Qt::CaseInsensitive) == 0 || ext.compare("alac", Qt::CaseInsensitive) == 0 ||
            ext.compare("tta", Qt::CaseInsensitive) == 0 || ext.compare("flac", Qt::CaseInsensitive) == 0 ||
            ext.compare("wma", Qt::CaseInsensitive) == 0 || ext.compare("wv", Qt::CaseInsensitive) == 0);
}

bool PlayList::isEmpty()
{
    if (fileList.isEmpty())
        return true;
    else
        return false;
}

void PlayList::add(const QString &fileName)
{
    if (fixSuffix(fileName))
    {
        if ((int)(m_player->getFileSecond(fileName)) >= lengthFilter)
        {
            fileList.append(fileName);
            timeList.append(m_player->getFileTotalTime(fileName));
        }
    }
    tableUpdate();
}

void PlayList::insert(int index, const QString &fileName)
{
    if (fixSuffix(fileName))
    {
        if ((int)(m_player->getFileSecond(fileName)) >= lengthFilter)
        {
            if (index < curIndex)
                ++curIndex;
            fileList.insert(index, fileName);
            timeList.insert(index, m_player->getFileTotalTime(fileName));
        }
    }
    tableUpdate();
}

void PlayList::remove(int index)
{
    if (index <= curIndex && index > -1)
        --curIndex;
    fileList.removeAt(index);
    timeList.removeAt(index);
    tableUpdate();
}

void PlayList::clearAll()
{
    fileList.clear();
    timeList.clear();
    curIndex = 0;
    tableUpdate();
}

int PlayList::getLength()
{
    return fileList.length();
}

int PlayList::getIndex()
{
    if (!fileList.isEmpty())
    {
        return curIndex;
    }
    else
    {
        return -1;
    }
}

QString PlayList::next(bool isLoop)
{
    if (!fileList.isEmpty())
    {
        if (isLoop)
        {
            if (curIndex < fileList.length() - 1)
            {
                ++curIndex;
            }
            else
            {
                curIndex = 0;
            }
            ui->playListTable->selectRow(curIndex);
            tableUpdate();
            return fileList[curIndex];
        }
        else
        {
            if (curIndex < fileList.length() - 1)
            {
                ++curIndex;
                ui->playListTable->selectRow(curIndex);
                tableUpdate();
                return fileList[curIndex];
            }
            else
            {
                return "stop";
            }
        }
    }
    return "";
}

QString PlayList::previous(bool isLoop)
{
    if (!fileList.isEmpty())
    {
        if (isLoop)
        {
            if (curIndex == 0)
            {
                curIndex = fileList.length() - 1;
            }
            else
            {
                --curIndex;
            }
            ui->playListTable->selectRow(curIndex);
            tableUpdate();
            return fileList[curIndex];
        }
        else
        {
            if (curIndex > 0)
            {
                --curIndex;
                ui->playListTable->selectRow(curIndex);
                tableUpdate();
                return fileList[curIndex];
            }
            else
            {
                return "stop";
            }
        }
    }
    return "";
}

QString PlayList::playIndex(int index)
{
    curIndex = index;
    ui->playListTable->selectRow(curIndex);
    tableUpdate();
    return fileList[curIndex];
}

QString PlayList::getFileNameForIndex(int index)
{
    return fileList[index];
}

QString PlayList::getCurFile()
{
    return fileList[curIndex];
}

QString PlayList::playLast()
{
    if (!fileList.isEmpty())
    {
        curIndex = fileList.length() - 1;
        ui->playListTable->selectRow(curIndex);
        tableUpdate();
        return fileList[curIndex];
    }
    return "";
}

void PlayList::tableUpdate()
{
    ui->playListTable->clear();
    ui->playListTable->setRowCount(getLength());
    int count = fileList.size();
    for (int i = 0; i < count; i++)
    {
        QString   fileName = fileList[i];
        QFileInfo fileInfo(fileName);

        QTableWidgetItem *item     = new QTableWidgetItem(fileInfo.fileName());
        QTableWidgetItem *timeItem = new QTableWidgetItem(timeList[i]);

        if (i == curIndex)
        {
            item->setBackgroundColor(QColor(128, 255, 0, 128));
            timeItem->setBackgroundColor(QColor(128, 255, 0, 128));
        }

        ui->playListTable->setItem(i, 0, item);
        ui->playListTable->setItem(i, 1, timeItem);
    }
}

void PlayList::on_deleteButton_clicked()
{
    QItemSelectionModel *selectionModel = ui->playListTable->selectionModel();
    QModelIndexList selected = selectionModel->selectedIndexes();
    QMap<int, int>       rowMap;
    foreach (QModelIndex index, selected)
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
    QList<QUrl> urls = event->mimeData()->urls();
    if (urls.isEmpty())
        return;

    int     urlCount = urls.size();
    QString fileName;

    for (int i = 0; i < urlCount; i++)
    {
        fileName = urls[i].toLocalFile();
        if (fixSuffix(fileName))
        {
            add(fileName);
        }
    }
}

void PlayList::on_playListTable_cellDoubleClicked(int row, int)
{
    curIndex = row;
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
    int count = fileNames.size();
    for (int i = 0; i < count; i++)
    {
        QString fileName = fileNames[i];
        add(fileName);
    }
}

void PlayList::on_searchButton_clicked()
{
    if (!fileList.isEmpty() && !ui->searchEdit->text().isEmpty())
    {
        int resultIndex = -1;
        int count       = fileList.size();
        for (int i = 0; i < count; i++)
        {
            QString   fileName = fileList[i];
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
    if (!fileList.isEmpty() && !ui->searchEdit->text().isEmpty())
    {
        int resultIndex = -1;
        int start       = ui->playListTable->currentRow() + 1;
        int count       = fileList.size();

        if (start < count)
            for (int i = start; i < count; i++)
            {
                QString   fileName = fileList[i];
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
                                   lengthFilter,
                                   0,
                                   2147483647,
                                   1,
                                   &ok);
    if (ok)
        lengthFilter = set;
}

void PlayList::saveToFile(const QString &fileName)
{
    QFile file(fileName);
    file.open(QIODevice::WriteOnly);
    QDataStream stream(&file);
    stream << (quint32)0x61727487 << fileList << timeList << curIndex;
    file.close();
}

void PlayList::readFromFile(const QString &fileName)
{
    QFile file(fileName);
    file.open(QIODevice::ReadOnly);
    QDataStream stream(&file);
    quint32     magic;
    stream >> magic;
    if (magic == 0x61727487)
    {
        stream >> fileList;
        stream >> timeList;
        stream >> curIndex;
    }
    file.close();
    tableUpdate();
}
