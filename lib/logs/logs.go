package logs

import (
"github.com/Sirupsen/logrus"
"os"

"github.com/logrusorgru/aurora"
)


var OsGetenv = os.Getenv

const colorsFlag bool = true

type Logger struct {
	internalLog *logrus.Logger
	prefix string
	constFileds map[string]interface{}
	CS aurora.Aurora
}


func NewLogger(name string) *Logger {
	obj := new (Logger)
	obj.CS = aurora.NewAurora(colorsFlag)
	obj.internalLog = logrus.New()
	obj.internalLog.Level = logrus.DebugLevel
	if ( os.Getenv("C3LOGFMT") == "JSON" ) {
		obj.internalLog.Formatter = &logrus.JSONFormatter{}
	}
	obj.prefix = name
	return obj
}

func (obj Logger)SubLogger (name string) *Logger{
	retval := NewLogger(obj.prefix + name)
	retval.constFileds = obj.constFileds

	return retval
}

func (obj Logger) Info(msg string, fields logrus.Fields){
	for key, value := range obj.constFileds {
		fields[key] = value
	}
	obj.internalLog.WithFields( fields ).Debug( obj.CS.Sprintf( obj.CS.Cyan("[ %s ] %s"), obj.CS.Magenta( obj.prefix ) , msg ) )
}

func (obj Logger)  Log(msg string ){
	obj.internalLog.Info( obj.CS.Sprintf("[ %s ] %s", obj.CS.Green( obj.prefix ) , msg )  )
}

func (obj Logger)  Debug(msg string ){
	obj.internalLog.Debug( obj.CS.Sprintf(obj.CS.Cyan ( "[ %s ] %s" ), obj.CS.Green( obj.prefix ) , msg ) )
}

func (obj Logger)  Error(err error ){
	obj.internalLog.Error( obj.CS.Sprintf( "[ %s ] %s" , obj.CS.Green( obj.prefix ) , obj.CS.Cyan( err ) ) )
}
