

#define protected public
#define private public

#include "moc.h"
#include "_cgo_export.h"

#include <QByteArray>
#include <QCamera>
#include <QCameraImageCapture>
#include <QChildEvent>
#include <QDBusPendingCall>
#include <QDBusPendingCallWatcher>
#include <QEvent>
#include <QExtensionFactory>
#include <QExtensionManager>
#include <QGraphicsObject>
#include <QGraphicsWidget>
#include <QLayout>
#include <QList>
#include <QMediaPlaylist>
#include <QMediaRecorder>
#include <QMetaMethod>
#include <QObject>
#include <QOffscreenSurface>
#include <QPaintDevice>
#include <QPaintDeviceWindow>
#include <QPdfWriter>
#include <QQuickItem>
#include <QRadioData>
#include <QSignalSpy>
#include <QString>
#include <QTime>
#include <QTimer>
#include <QTimerEvent>
#include <QWidget>
#include <QWindow>


class ClientQT: public QObject
{
Q_OBJECT
public:
	ClientQT(QObject *parent = Q_NULLPTR) : QObject(parent) {qRegisterMetaType<quintptr>("quintptr");ClientQT_ClientQT_QRegisterMetaType();ClientQT_ClientQT_QRegisterMetaTypes();callbackClientQT_Constructor(this);};
	void Signal_ReloadUI() { callbackClientQT_ReloadUI(this); };
	 ~ClientQT() { callbackClientQT_DestroyClientQT(this); };
	bool event(QEvent * e) { return callbackClientQT_Event(this, e) != 0; };
	bool eventFilter(QObject * watched, QEvent * event) { return callbackClientQT_EventFilter(this, watched, event) != 0; };
	void childEvent(QChildEvent * event) { callbackClientQT_ChildEvent(this, event); };
	void connectNotify(const QMetaMethod & sign) { callbackClientQT_ConnectNotify(this, const_cast<QMetaMethod*>(&sign)); };
	void customEvent(QEvent * event) { callbackClientQT_CustomEvent(this, event); };
	void deleteLater() { callbackClientQT_DeleteLater(this); };
	void Signal_Destroyed(QObject * obj) { callbackClientQT_Destroyed(this, obj); };
	void disconnectNotify(const QMetaMethod & sign) { callbackClientQT_DisconnectNotify(this, const_cast<QMetaMethod*>(&sign)); };
	void Signal_ObjectNameChanged(const QString & objectName) { QByteArray taa2c4f = objectName.toUtf8(); Moc_PackedString objectNamePacked = { const_cast<char*>(taa2c4f.prepend("WHITESPACE").constData()+10), taa2c4f.size()-10 };callbackClientQT_ObjectNameChanged(this, objectNamePacked); };
	void timerEvent(QTimerEvent * event) { callbackClientQT_TimerEvent(this, event); };
	
signals:
	void reloadUI();
public slots:
private:
};

Q_DECLARE_METATYPE(ClientQT*)


void ClientQT_ClientQT_QRegisterMetaTypes() {
}

void ClientQT_ConnectReloadUI(void* ptr)
{
	QObject::connect(static_cast<ClientQT*>(ptr), static_cast<void (ClientQT::*)()>(&ClientQT::reloadUI), static_cast<ClientQT*>(ptr), static_cast<void (ClientQT::*)()>(&ClientQT::Signal_ReloadUI));
}

void ClientQT_DisconnectReloadUI(void* ptr)
{
	QObject::disconnect(static_cast<ClientQT*>(ptr), static_cast<void (ClientQT::*)()>(&ClientQT::reloadUI), static_cast<ClientQT*>(ptr), static_cast<void (ClientQT::*)()>(&ClientQT::Signal_ReloadUI));
}

void ClientQT_ReloadUI(void* ptr)
{
	static_cast<ClientQT*>(ptr)->reloadUI();
}

int ClientQT_ClientQT_QRegisterMetaType()
{
	return qRegisterMetaType<ClientQT*>();
}

int ClientQT_ClientQT_QRegisterMetaType2(char* typeName)
{
	return qRegisterMetaType<ClientQT*>(const_cast<const char*>(typeName));
}

int ClientQT_ClientQT_QmlRegisterType()
{
#ifdef QT_QML_LIB
	return qmlRegisterType<ClientQT>();
#else
	return 0;
#endif
}

int ClientQT_ClientQT_QmlRegisterType2(char* uri, int versionMajor, int versionMinor, char* qmlName)
{
#ifdef QT_QML_LIB
	return qmlRegisterType<ClientQT>(const_cast<const char*>(uri), versionMajor, versionMinor, const_cast<const char*>(qmlName));
#else
	return 0;
#endif
}

void* ClientQT___dynamicPropertyNames_atList(void* ptr, int i)
{
	return new QByteArray(static_cast<QList<QByteArray>*>(ptr)->at(i));
}

void ClientQT___dynamicPropertyNames_setList(void* ptr, void* i)
{
	static_cast<QList<QByteArray>*>(ptr)->append(*static_cast<QByteArray*>(i));
}

void* ClientQT___dynamicPropertyNames_newList(void* ptr)
{
	Q_UNUSED(ptr);
	return new QList<QByteArray>;
}

void* ClientQT___findChildren_atList2(void* ptr, int i)
{
	return const_cast<QObject*>(static_cast<QList<QObject*>*>(ptr)->at(i));
}

void ClientQT___findChildren_setList2(void* ptr, void* i)
{
	static_cast<QList<QObject*>*>(ptr)->append(static_cast<QObject*>(i));
}

void* ClientQT___findChildren_newList2(void* ptr)
{
	Q_UNUSED(ptr);
	return new QList<QObject*>;
}

void* ClientQT___findChildren_atList3(void* ptr, int i)
{
	return const_cast<QObject*>(static_cast<QList<QObject*>*>(ptr)->at(i));
}

void ClientQT___findChildren_setList3(void* ptr, void* i)
{
	static_cast<QList<QObject*>*>(ptr)->append(static_cast<QObject*>(i));
}

void* ClientQT___findChildren_newList3(void* ptr)
{
	Q_UNUSED(ptr);
	return new QList<QObject*>;
}

void* ClientQT___findChildren_atList(void* ptr, int i)
{
	return const_cast<QObject*>(static_cast<QList<QObject*>*>(ptr)->at(i));
}

void ClientQT___findChildren_setList(void* ptr, void* i)
{
	static_cast<QList<QObject*>*>(ptr)->append(static_cast<QObject*>(i));
}

void* ClientQT___findChildren_newList(void* ptr)
{
	Q_UNUSED(ptr);
	return new QList<QObject*>;
}

void* ClientQT___children_atList(void* ptr, int i)
{
	return const_cast<QObject*>(static_cast<QList<QObject *>*>(ptr)->at(i));
}

void ClientQT___children_setList(void* ptr, void* i)
{
	static_cast<QList<QObject *>*>(ptr)->append(static_cast<QObject*>(i));
}

void* ClientQT___children_newList(void* ptr)
{
	Q_UNUSED(ptr);
	return new QList<QObject *>;
}

void* ClientQT_NewClientQT(void* parent)
{
	if (dynamic_cast<QCameraImageCapture*>(static_cast<QObject*>(parent))) {
		return new ClientQT(static_cast<QCameraImageCapture*>(parent));
	} else if (dynamic_cast<QDBusPendingCallWatcher*>(static_cast<QObject*>(parent))) {
		return new ClientQT(static_cast<QDBusPendingCallWatcher*>(parent));
	} else if (dynamic_cast<QExtensionFactory*>(static_cast<QObject*>(parent))) {
		return new ClientQT(static_cast<QExtensionFactory*>(parent));
	} else if (dynamic_cast<QExtensionManager*>(static_cast<QObject*>(parent))) {
		return new ClientQT(static_cast<QExtensionManager*>(parent));
	} else if (dynamic_cast<QGraphicsObject*>(static_cast<QObject*>(parent))) {
		return new ClientQT(static_cast<QGraphicsObject*>(parent));
	} else if (dynamic_cast<QGraphicsWidget*>(static_cast<QObject*>(parent))) {
		return new ClientQT(static_cast<QGraphicsWidget*>(parent));
	} else if (dynamic_cast<QLayout*>(static_cast<QObject*>(parent))) {
		return new ClientQT(static_cast<QLayout*>(parent));
	} else if (dynamic_cast<QMediaPlaylist*>(static_cast<QObject*>(parent))) {
		return new ClientQT(static_cast<QMediaPlaylist*>(parent));
	} else if (dynamic_cast<QMediaRecorder*>(static_cast<QObject*>(parent))) {
		return new ClientQT(static_cast<QMediaRecorder*>(parent));
	} else if (dynamic_cast<QOffscreenSurface*>(static_cast<QObject*>(parent))) {
		return new ClientQT(static_cast<QOffscreenSurface*>(parent));
	} else if (dynamic_cast<QPaintDeviceWindow*>(static_cast<QObject*>(parent))) {
		return new ClientQT(static_cast<QPaintDeviceWindow*>(parent));
	} else if (dynamic_cast<QPdfWriter*>(static_cast<QObject*>(parent))) {
		return new ClientQT(static_cast<QPdfWriter*>(parent));
	} else if (dynamic_cast<QQuickItem*>(static_cast<QObject*>(parent))) {
		return new ClientQT(static_cast<QQuickItem*>(parent));
	} else if (dynamic_cast<QRadioData*>(static_cast<QObject*>(parent))) {
		return new ClientQT(static_cast<QRadioData*>(parent));
	} else if (dynamic_cast<QSignalSpy*>(static_cast<QObject*>(parent))) {
		return new ClientQT(static_cast<QSignalSpy*>(parent));
	} else if (dynamic_cast<QWidget*>(static_cast<QObject*>(parent))) {
		return new ClientQT(static_cast<QWidget*>(parent));
	} else if (dynamic_cast<QWindow*>(static_cast<QObject*>(parent))) {
		return new ClientQT(static_cast<QWindow*>(parent));
	} else {
		return new ClientQT(static_cast<QObject*>(parent));
	}
}

void ClientQT_DestroyClientQT(void* ptr)
{
	static_cast<ClientQT*>(ptr)->~ClientQT();
}

void ClientQT_DestroyClientQTDefault(void* ptr)
{
	Q_UNUSED(ptr);

}

char ClientQT_EventDefault(void* ptr, void* e)
{
	return static_cast<ClientQT*>(ptr)->QObject::event(static_cast<QEvent*>(e));
}

char ClientQT_EventFilterDefault(void* ptr, void* watched, void* event)
{
	return static_cast<ClientQT*>(ptr)->QObject::eventFilter(static_cast<QObject*>(watched), static_cast<QEvent*>(event));
}

void ClientQT_ChildEventDefault(void* ptr, void* event)
{
	static_cast<ClientQT*>(ptr)->QObject::childEvent(static_cast<QChildEvent*>(event));
}

void ClientQT_ConnectNotifyDefault(void* ptr, void* sign)
{
	static_cast<ClientQT*>(ptr)->QObject::connectNotify(*static_cast<QMetaMethod*>(sign));
}

void ClientQT_CustomEventDefault(void* ptr, void* event)
{
	static_cast<ClientQT*>(ptr)->QObject::customEvent(static_cast<QEvent*>(event));
}

void ClientQT_DeleteLaterDefault(void* ptr)
{
	static_cast<ClientQT*>(ptr)->QObject::deleteLater();
}

void ClientQT_DisconnectNotifyDefault(void* ptr, void* sign)
{
	static_cast<ClientQT*>(ptr)->QObject::disconnectNotify(*static_cast<QMetaMethod*>(sign));
}

void ClientQT_TimerEventDefault(void* ptr, void* event)
{
	static_cast<ClientQT*>(ptr)->QObject::timerEvent(static_cast<QTimerEvent*>(event));
}



#include "moc_moc.h"
