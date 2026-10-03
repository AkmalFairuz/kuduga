package com.kuduga.app

import android.annotation.SuppressLint
import android.bluetooth.BluetoothAdapter
import android.bluetooth.BluetoothManager
import android.content.Context
import android.content.Intent
import android.graphics.Bitmap
import com.dantsu.escposprinter.EscPosPrinter
import com.dantsu.escposprinter.EscPosPrinterCommands
import com.dantsu.escposprinter.connection.DeviceConnection
import com.dantsu.escposprinter.connection.bluetooth.BluetoothConnection
import com.dantsu.escposprinter.connection.bluetooth.BluetoothPrintersConnections
import com.dantsu.escposprinter.textparser.PrinterTextParserImg


class ThermalPrinter {

    companion object {
        const val REQUEST_BLUETOOTH = 104
        const val REQUEST_BLUETOOTH_SETTINGS = 105

        @SuppressLint("ServiceCast", "MissingPermission")
        fun requestTurnOnBluetooth(activity: MainActivity, context: Context) {
            val bluetoothManager = context.getSystemService(Context.BLUETOOTH_SERVICE) as BluetoothManager
            if(!bluetoothManager.adapter.isEnabled) {
                val enableBtIntent = Intent(BluetoothAdapter.ACTION_REQUEST_ENABLE)
                activity.startActivityForResult(enableBtIntent, REQUEST_BLUETOOTH);
            }
        }

        fun isBluetoothEnabled(context: Context): Boolean {
            val bluetoothManager = context.getSystemService(Context.BLUETOOTH_SERVICE) as BluetoothManager
            return bluetoothManager.adapter.isEnabled
        }

    }

    val conns: BluetoothPrintersConnections = BluetoothPrintersConnections()

    fun getBluetoothPrinters(): Array<out BluetoothConnection>? {
        return conns.list
    }

    fun getBluetoothPrinterByAddress(address: String): BluetoothConnection? {
        for(conn in conns.list!!) {
            if(conn.device.address == address) {
                return conn
            }
        }
        return null
    }

    fun printFormatted(conn: DeviceConnection, text: String, width: Float, dpi: Int) {
        conn.connect()
        val printer = EscPosPrinter(conn, dpi, width, 32)
        try {
            printer.printFormattedTextAndCut(text, 1)
        }finally {
            printer.disconnectPrinter()
            conn.disconnect()
        }
    }
}
