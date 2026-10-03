package com.kuduga.app

import android.annotation.SuppressLint
import android.app.Activity
import android.app.NotificationChannel
import android.app.NotificationManager
import android.content.Intent
import android.os.Build
import android.provider.Settings
import com.dantsu.escposprinter.textparser.PrinterTextParserImg
import io.flutter.embedding.android.FlutterFragmentActivity
import io.flutter.embedding.engine.FlutterEngine
import io.flutter.plugin.common.MethodChannel
import io.flutter.plugins.GeneratedPluginRegistrant
import java.lang.Exception
import java.util.UUID

class MainActivity: FlutterFragmentActivity() {

    private val DATA_CHANNEL = "finance.app.channel"

    private val thermalPrinters = mutableMapOf<String, ThermalPrinter>()
    private var pendingBluetoothResult: MethodChannel.Result? = null
    private var pendingBluetoothSettingsResult: MethodChannel.Result? = null

    @SuppressLint("MissingPermission")
    override fun configureFlutterEngine(flutterEngine: FlutterEngine) {
        GeneratedPluginRegistrant.registerWith(flutterEngine)
        MethodChannel(flutterEngine.dartExecutor.binaryMessenger, DATA_CHANNEL).setMethodCallHandler { call, result ->
            val args = call.arguments as? Map<*, *>
            var handled = false
            when(call.method) {
                "createNotificationChannel" -> {
                    if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O && args != null) {
                        // Create the NotificationChannel.
                        val name = args["name"] as String
                        val importance = NotificationManager.IMPORTANCE_HIGH
                        val mChannel = NotificationChannel(args["id"] as String, name, importance)
                        val notificationManager = getSystemService(NOTIFICATION_SERVICE) as NotificationManager
                        notificationManager.createNotificationChannel(mChannel)
                        result.success(true)
                    }else {
                        result.error("notificationChannelNotCreated", "", null)
                    }
                    handled = true
                }
                "createThermalPrinter" -> {
                    val thermalPrinter = ThermalPrinter()
                    val uuid = UUID.randomUUID().toString()
                    thermalPrinters[uuid] = thermalPrinter
                    result.success(uuid)
                    handled = true
                }
                "requestBluetoothSettings" -> {
                    if(pendingBluetoothSettingsResult != null) {
                        result.success(false)
                    } else {
                        val bluetoothIntent = Intent(Settings.ACTION_BLUETOOTH_SETTINGS)
                        startActivityForResult(
                            bluetoothIntent,
                            ThermalPrinter.REQUEST_BLUETOOTH_SETTINGS
                        )
                        pendingBluetoothSettingsResult = result
                    }
                    handled = true
                }
                "destroyThermalPrinter" -> {
                    if(args != null) {
                        thermalPrinters.remove(args["uuid"])
                        result.success("success")
                    }else{
                        result.error("printerNotFound", "Printer not found", null)
                    }
                    handled = true
                }
                "isBluetoothEnabled" -> {
                    result.success(ThermalPrinter.isBluetoothEnabled(getApplicationContext()))
                    handled = true
                }
                "requestBluetooth" -> {
                    if(pendingBluetoothResult != null) {
                        result.success(false)
                    }else {
                        pendingBluetoothResult = result
                        ThermalPrinter.requestTurnOnBluetooth(this, getApplicationContext());
                    }
                    handled = true;
                }
                "invokeThermalPrinter" -> {
                    if(args != null) {
                        val uuid = args["uuid"]
                        val method = args["method"]
                        val thermalPrinter = thermalPrinters[uuid]
                        if(thermalPrinter != null) {
                            when (method) {
                                "getPrinters" -> {
                                    val printers = thermalPrinter.getBluetoothPrinters()
                                    val response = arrayListOf<HashMap<String, String>>();
                                    if (printers != null) {
                                        for (printer in printers) {
                                            val map = HashMap<String, String>()
                                            map["name"] = printer.device.name
                                            map["address"] = printer.device.address
                                            response.add(map);
                                        }
                                    }
                                    result.success(response)
                                    handled = true
                                }
                                "print" -> {
                                    val connAddress = args["address"] as String
                                    val text = args["text"] as String
                                    val dpi = args["dpi"] as Int
                                    val width = args["width"] as Int
                                    val conn = thermalPrinter.getBluetoothPrinterByAddress(connAddress)
                                    if (conn != null) {
                                        try {
                                            thermalPrinter.printFormatted(conn, text, width.toFloat(), dpi)
                                            result.success("success")
                                        } catch (e: Exception) {
                                            result.error("printError", e.javaClass.name + ": " + e.message, null)
                                        }
                                    } else {
                                        result.error("printerIsNull", "Printer is not found", null)
                                    }
                                    handled = true
                                }
                            }
                        }
                    }
                }
            }
            if(!handled) {
                result.notImplemented()
            }
        }
    }

    override fun onActivityResult(requestCode: Int, resultCode: Int, data: Intent?) {
        if(requestCode == ThermalPrinter.REQUEST_BLUETOOTH) {
            pendingBluetoothResult?.success(resultCode == Activity.RESULT_OK)
            pendingBluetoothResult = null
        } else if (requestCode == ThermalPrinter.REQUEST_BLUETOOTH_SETTINGS) {
            pendingBluetoothSettingsResult?.success(resultCode == Activity.RESULT_OK)
            pendingBluetoothSettingsResult = null
        }
        super.onActivityResult(requestCode, resultCode, data)
    }
}
