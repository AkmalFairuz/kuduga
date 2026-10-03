import 'package:finance_app/model/thermal_printer.dart';
import 'package:finance_app/service/printer/printable.dart';
import 'package:finance_app/service/thermal_printer_service.dart';
import 'package:finance_app/utils/alert.dart';
import 'package:finance_app/utils/screens.dart';
import 'package:finance_app/widgets/input/input.dart';
import 'package:finance_app/widgets/layout/line.dart';
import 'package:finance_app/widgets/select/select_box.dart';
import 'package:finance_app/widgets/text/text_link.dart';
import 'package:flutter/material.dart';

import '../../widgets/text/text_app_bar.dart';

class SelectPrinterScreen extends StatefulWidget {
  const SelectPrinterScreen({Key? key, required this.printable})
      : super(key: key);

  final ThermalPrintable printable;

  @override
  State<SelectPrinterScreen> createState() => _SelectPrinterState();
}

class _SelectPrinterState extends State<SelectPrinterScreen> {
  final ThermalPrinterService _service = ThermalPrinterService();
  bool _scanned = false;
  List<ThermalPrinterConnection> _printers = [];
  final printerWidthController = SelectController<int?>(null);
  final headerController = TextEditingController();
  final footerController = TextEditingController();

  @override
  void initState() {
    super.initState();

    _initialize();
  }

  void _initialize() async {
    await _service.initialize();

    printerWidthController.setValue(_service.savedWidthSize);
    headerController.text = _service.savedHeaderText ?? "";
    footerController.text = _service.savedFooterText ?? "";

    if (!await _service.checkPermissions()) {
      Alert.message("Izin bluetooth diperlukan")
          .then((value) => Screens.back());
      return;
    }

    var isBluetoothOn = await _service.isBluetoothEnabled();
    if (!isBluetoothOn) {
      if (!await _service.requestEnableBluetooth()) {
        Alert.message("Bluetooth tidak aktif").then((value) => Screens.back());
        return;
      }
    }

    _scanDevices();
  }

  void _scanDevices() async {
    setState(() {
      _scanned = false;
    });

    var printers = await _service.getThermalPrinters();
    if (_service.savedThermalPrinterUsed != null) {
      for (final printer in printers) {
        if (printer.toString() == _service.savedThermalPrinterUsed) {
          printers.remove(printer);
          printers = [printer, ...printers];
          break;
        }
      }
    }

    setState(() {
      _printers = printers;
      _scanned = true;
    });
  }

  void _onAddPrinter() async {
    await Alert.show(
        closeMessage: "Lanjutkan",
        title: "Tambahkan Thermal Printer",
        message:
            "Jika anda tidak menemukan printer anda, anda harus menambahkannya di Setelan Bluetooth.\n\nBuka Bluetooth di Setelan/Pengaturan perangkat anda. Kemudian tambahkan thermal printer anda, biasanya anda akan diminta untuk memasukkan PIN thermal printer anda (biasanya PIN 0000 atau 9999). Setelah terhubung dengan thermal printer anda, anda harus memilih thermal printer anda di aplikasi");
    await _service.requestBluetoothSettings();
    _scanDevices();
  }

  @override
  void dispose() {
    _service.destroy();

    super.dispose();
  }

  Future<void> _onRequestPrint(ThermalPrinterConnection connection) async {
    if (printerWidthController.value() == null) {
      Alert.message("Panjang kertas harus dipilih");
      return;
    }
    Alert.withLoading((done) async {
      await _service.printFormatted(connection, widget.printable,
          paperWidth: printerWidthController.value()!,
          headerText: headerController.text,
          footerText: footerController.text);
      done();
    }, message: "Printing...");
  }

  Widget _buildContent() {
    return _printers.isNotEmpty
        ? ListView.builder(
            itemBuilder: (context, index) {
              var printer = _printers[index];
              return ListTile(
                leading: const Icon(Icons.bluetooth, size: 36),
                title: Text(printer.name),
                subtitle: Text(printer.address),
                trailing: _service.savedThermalPrinterUsed == printer.toString()
                    ? const Text(
                        "Terakhir digunakan",
                        style: TextStyle(
                            fontWeight: FontWeight.bold, fontSize: 12),
                      )
                    : null,
                onTap: () {
                  _onRequestPrint(printer);
                },
              );
            },
            itemCount: _printers.length)
        : Center(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.center,
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                const Text('Tidak ada printer yang tersedia'),
                TextLink(text: 'Tambahkan Printer', onPressed: _onAddPrinter)
              ],
            ),
          );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      resizeToAvoidBottomInset: false,
      appBar: AppBar(
        title: const AppBarTitle('Pilih Printer'),
        actions: _scanned
            ? [
                IconButton(
                    onPressed: _onAddPrinter, icon: const Icon(Icons.add)),
                IconButton(
                    onPressed: () {
                      _scanDevices();
                    },
                    icon: const Icon(Icons.refresh)),
              ]
            : null,
      ),
      body: Column(children: [
        const SizedBox(height: 16),
        Padding(
            padding: const EdgeInsets.symmetric(horizontal: 16),
            child: Row(
              children: [
                Expanded(
                    child: SelectBox(
                        label: "Panjang Kertas",
                        controller: printerWidthController,
                        items: ThermalPrinterService.supportedPaperWidth
                            .map((e) => SelectItem(e, "${e}mm"))
                            .toList()))
              ],
            )),
        const SizedBox(height: 8),
        Padding(
          padding: const EdgeInsets.symmetric(horizontal: 16),
          child: Input(
            controller: headerController,
            minLines: 2,
            maxLines: 3,
            keyboardType: TextInputType.multiline,
            padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 14),
            label: "Teks bagian atas",
          ),
        ),
        const SizedBox(height: 8),
        Padding(
          padding: const EdgeInsets.symmetric(horizontal: 16),
          child: Input(
            controller: footerController,
            minLines: 2,
            maxLines: 3,
            keyboardType: TextInputType.multiline,
            padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 14),
            label: "Teks bagian bawah",
          ),
        ),
        const SizedBox(height: 16),
        const Line(),
        Expanded(
            child: _scanned
                ? _buildContent()
                : const Center(child: CircularProgressIndicator()))
      ]),
    );
  }
}
