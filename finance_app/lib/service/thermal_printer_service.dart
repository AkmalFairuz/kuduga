import 'package:finance_app/model/thermal_printer.dart';
import 'package:finance_app/service/printer/printable.dart';
import 'package:finance_app/utils/native_bindings.dart';
import 'package:permission_handler/permission_handler.dart';

import '../utils/local_storage.dart';

class ThermalPrinterService {
  // max characters per line on normal font size.
  static final maxCharactersPerLine = {
    // 58mm = 32 or 33
    58: 32,
    // 80mm = 48 to 52
    80: 48
  };

  static final supportedPaperWidth = [58, 80];

  ThermalPrinterService();

  final _channel = NativeBindings.channel;
  String? uuid;
  int? savedWidthSize;
  String? savedThermalPrinterUsed;
  String? savedHeaderText;
  String? savedFooterText;

  Future<void> initialize() async {
    if (uuid != null) {
      return;
    }

    savedWidthSize = await LocalStorage.get("thermalPrinterWidthSize");
    savedThermalPrinterUsed = await LocalStorage.get("lastThermalPrinterUsed");
    savedHeaderText =
        await LocalStorage.get("thermalPrinterHeaderText", def: "");
    savedFooterText =
        await LocalStorage.get("thermalPrinterFooterText", def: "");

    return _channel.invokeMethod("createThermalPrinter").then((value) {
      uuid = value as String;
    });
  }

  Future<void> requestBluetoothSettings() {
    return _channel.invokeMethod("requestBluetoothSettings");
  }

  Future<bool> checkPermissions() async {
    if (!await Permission.bluetooth.request().isGranted) {
      return false;
    }
    if (!await Permission.bluetoothScan.request().isGranted) {
      return false;
    }
    if (!await Permission.bluetoothConnect.request().isGranted) {
      return false;
    }
    return true;
  }

  Future<bool> isBluetoothEnabled() {
    return _channel
        .invokeMethod("isBluetoothEnabled")
        .then((value) => value as bool);
  }

  Future<bool> requestEnableBluetooth() {
    return _channel
        .invokeMethod("requestBluetooth")
        .then((value) => value as bool);
  }

  Future<bool> requestPermission() {
    return _channel
        .invokeMethod("requestThermalPrinterPermission")
        .then((value) => value as bool);
  }

  Future<void> destroy() {
    return _channel
        .invokeMethod("destroyThermalPrinter", {'uuid': uuid}).then((value) {});
  }

  Future<List<ThermalPrinterConnection>> getThermalPrinters() {
    return _interact("getPrinters", null)
        .then((connections) => (connections as List<Object?>).map((conn) {
              var conn2 = conn as Map<Object?, Object?>;
              return ThermalPrinterConnection(
                  conn2["name"] as String, conn2["address"] as String);
            }).toList());
  }

  Future<void> saveThermalPrinterInfo(ThermalPrinterConnection conn,
      int paperWidth, String headerText, String footerText) async {
    savedWidthSize = paperWidth;
    savedThermalPrinterUsed = conn.toString();
    await LocalStorage.set("thermalPrinterWidthSize", paperWidth);
    await LocalStorage.set("lastThermalPrinterUsed", conn.toString());
    await LocalStorage.set("thermalPrinterHeaderText", headerText);
    await LocalStorage.set("thermalPrinterFooterText", footerText);
  }

  String _processTags(String text) {
    text = text.replaceAll("[", "");
    text = text.replaceAll("]", "");
    text = text.replaceAll("<", "");
    text = text.replaceAll(">", "");
    text = text.replaceAll("@big@", "<font size='big'>");
    text = text.replaceAll("@big2@", "<font size='big-2'>");
    text = text.replaceAll("@big3@", "<font size='big-3'>");
    text = text.replaceAll("@wide@", "<font size='wide'>");
    text = text.replaceAll("@tall@", "<font size='tall'>");
    text = text.replaceAll("@b@", "<b>");
    text = text.replaceAll("@/big@", "</font>");
    text = text.replaceAll("@/big2@", "</font>");
    text = text.replaceAll("@/big3@", "</font>");
    text = text.replaceAll("@/wide@", "</font>");
    text = text.replaceAll("@/tall@", "</font>");
    text = text.replaceAll("@/b@", "</b>");
    return "[C]${text.split("\n").join("\n[C]")}";
  }

  Future<void> printFormatted(
      ThermalPrinterConnection connection, ThermalPrintable printable,
      {int paperWidth = 58,
      int printDpi = 203,
      String headerText = "",
      String footerText = ""}) {
    saveThermalPrinterInfo(connection, paperWidth, headerText, footerText);

    var rawText = "\n";
    if (headerText.isNotEmpty) {
      headerText = _processTags(headerText);
      rawText += "$headerText\n\n";
    }
    rawText += printable.thermalPrint(maxCharactersPerLine[paperWidth] ??
        32); // TODO: support other paper width
    if (footerText.isNotEmpty) {
      footerText = _processTags(footerText);
      rawText += "\n\n$footerText";
    }
    rawText += "\n";

    return _interact("print", {
      "text": rawText,
      "address": connection.address,
      "width": paperWidth,
      "dpi": printDpi,
    });
  }

  Future _interact<T>(String method, Map<String, dynamic>? args) {
    return _channel.invokeMethod(
        "invokeThermalPrinter", {'method': method, 'uuid': uuid, ...?args});
  }
}
