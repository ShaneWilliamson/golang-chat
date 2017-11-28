

#pragma once

#ifndef GO_MOC_H
#define GO_MOC_H

#include <stdint.h>

#ifdef __cplusplus
class ClientQT;
void ClientQT_ClientQT_QRegisterMetaTypes();
extern "C" {
#endif

struct Moc_PackedString { char* data; long long len; };
struct Moc_PackedList { void* data; long long len; };
void ClientQT_ConnectReloadUI(void* ptr);
void ClientQT_DisconnectReloadUI(void* ptr);
void ClientQT_ReloadUI(void* ptr);
int ClientQT_ClientQT_QRegisterMetaType();
int ClientQT_ClientQT_QRegisterMetaType2(char* typeName);
int ClientQT_ClientQT_QmlRegisterType();
int ClientQT_ClientQT_QmlRegisterType2(char* uri, int versionMajor, int versionMinor, char* qmlName);
void* ClientQT___dynamicPropertyNames_atList(void* ptr, int i);
void ClientQT___dynamicPropertyNames_setList(void* ptr, void* i);
void* ClientQT___dynamicPropertyNames_newList(void* ptr);
void* ClientQT___findChildren_atList2(void* ptr, int i);
void ClientQT___findChildren_setList2(void* ptr, void* i);
void* ClientQT___findChildren_newList2(void* ptr);
void* ClientQT___findChildren_atList3(void* ptr, int i);
void ClientQT___findChildren_setList3(void* ptr, void* i);
void* ClientQT___findChildren_newList3(void* ptr);
void* ClientQT___findChildren_atList(void* ptr, int i);
void ClientQT___findChildren_setList(void* ptr, void* i);
void* ClientQT___findChildren_newList(void* ptr);
void* ClientQT___children_atList(void* ptr, int i);
void ClientQT___children_setList(void* ptr, void* i);
void* ClientQT___children_newList(void* ptr);
void* ClientQT_NewClientQT(void* parent);
void ClientQT_DestroyClientQT(void* ptr);
void ClientQT_DestroyClientQTDefault(void* ptr);
char ClientQT_EventDefault(void* ptr, void* e);
char ClientQT_EventFilterDefault(void* ptr, void* watched, void* event);
void ClientQT_ChildEventDefault(void* ptr, void* event);
void ClientQT_ConnectNotifyDefault(void* ptr, void* sign);
void ClientQT_CustomEventDefault(void* ptr, void* event);
void ClientQT_DeleteLaterDefault(void* ptr);
void ClientQT_DisconnectNotifyDefault(void* ptr, void* sign);
void ClientQT_TimerEventDefault(void* ptr, void* event);
;

#ifdef __cplusplus
}
#endif

#endif