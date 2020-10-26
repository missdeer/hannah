#import <QUrl>
#import <QtCore>

#import "service.h"

#import <Cocoa/Cocoa.h>

#import "serviceslots.h"

@implementation HannahService

- (void)search:(NSPasteboard *)pboard userData:(NSString *)userData error:(NSString **)error
{
    if (![pboard canReadObjectForClasses:@[ [NSString class] ] options:@{}])
    {
        *error = NSLocalizedString(@"Error: the pasteboard doesn't contain a string.", nil);
        return;
    }

    // Try to parse a selector from the pasteboard contents.
    NSString *pasteboardString = [pboard stringForType:NSPasteboardTypeString];

    QString s = QString::fromNSString(pasteboardString);
    qDebug() << "search keyword:" << s;

    serviceSearch(s);
}

- (void)openUrl:(NSPasteboard *)pboard userData:(NSString *)userData error:(NSString **)error
{
    if (![pboard canReadObjectForClasses:@[ [NSString class] ] options:@{}])
    {
        *error = NSLocalizedString(@"Error: the pasteboard doesn't contain a string.", nil);
        return;
    }

    // Try to parse a selector from the pasteboard contents.
    NSString *pasteboardString = [pboard stringForType:NSPasteboardTypeString];

    QString s = QString::fromNSString(pasteboardString);
    qDebug() << "open Url:" << s;

    serviceOpenUrl(s);
}

- (void)openLink:(NSPasteboard *)pboard userData:(NSString *)userData error:(NSString **)error
{
    if (![pboard canReadObjectForClasses:@[ [NSString class] ] options:@{}])
    {
        *error = NSLocalizedString(@"Error: the pasteboard doesn't contain a string.", nil);
        return;
    }

    // Try to parse a selector from the pasteboard contents.
    NSString *pasteboardString = [pboard stringForType:NSPasteboardTypeString];

    QString s = QString::fromNSString(pasteboardString);
    qDebug() << "open link:" << s;

    serviceOpenLink(s);
}

- (void)appendToPlaylist:(NSPasteboard *)pboard userData:(NSString *)userData error:(NSString **)error
{
    if (![[pboard types] containsObject:NSFilenamesPboardType])
    {
        *error = NSLocalizedString(@"Error: the pasteboard doesn't contain directory(s)/file(s).", nil);
        return;
    }

    QStringList ss;
    NSArray *   fileArray = [pboard propertyListForType:NSFilenamesPboardType];
    for (NSString *filePath in fileArray)
    {
        ss << QString::fromNSString(filePath);
    }

    qDebug() << "appendToPlaylist:" << ss;
    serviceAppendToPlaylist(ss);
}

- (void)clearAndAddToPlaylist:(NSPasteboard *)pboard userData:(NSString *)userData error:(NSString **)error
{
    if (![[pboard types] containsObject:NSFilenamesPboardType])
    {
        *error = NSLocalizedString(@"Error: the pasteboard doesn't contain directory(s)/file(s).", nil);
        return;
    }

    QStringList ss;
    NSArray *   fileArray = [pboard propertyListForType:NSFilenamesPboardType];
    for (NSString *filePath in fileArray)
    {
        ss << QString::fromNSString(filePath);
    }

    qDebug() << "clearAndAddToPlaylist:" << ss;
    serviceClearAndAddToPlaylist(ss);
}

- (void)appendToPlaylistFile:(NSPasteboard *)pboard userData:(NSString *)userData error:(NSString **)error
{
    if (![[pboard types] containsObject:NSFilenamesPboardType])
    {
        *error = NSLocalizedString(@"Error: the pasteboard doesn't contain directory(s)/file(s).", nil);
        return;
    }

    QStringList ss;
    NSArray *   fileArray = [pboard propertyListForType:NSFilenamesPboardType];
    for (NSString *filePath in fileArray)
    {
        ss << QString::fromNSString(filePath);
    }

    qDebug() << "appendToPlaylistFile:" << ss;
    serviceAppendToPlaylistFile(ss);
}

- (void)clearAndAddToPlaylistFile:(NSPasteboard *)pboard userData:(NSString *)userData error:(NSString **)error
{
    if (![[pboard types] containsObject:NSFilenamesPboardType])
    {
        *error = NSLocalizedString(@"Error: the pasteboard doesn't contain directory(s)/file(s).", nil);
        return;
    }

    QStringList ss;
    NSArray *   fileArray = [pboard propertyListForType:NSFilenamesPboardType];
    for (NSString *filePath in fileArray)
    {
        ss << QString::fromNSString(filePath);
    }

    qDebug() << "clearAndAddToPlaylistFile:" << ss;
    serviceClearAndAddToPlaylistFile(ss);
}

@end

void registerHannahService()
{
    NSRegisterServicesProvider([HannahService new], @"Hannah");
    NSUpdateDynamicServices();
}
