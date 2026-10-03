import 'package:finance_app/screens/balance/deposit/deposit_content.dart';
import 'package:finance_app/utils/meta.dart';
import 'package:finance_app/utils/toast.dart';
import 'package:finance_app/widgets/input/input.dart';
import 'package:flutter/material.dart';
import 'package:image_gallery_saver/image_gallery_saver.dart';
import 'package:qr_flutter/qr_flutter.dart';
import 'package:screenshot/screenshot.dart';

class QrisContent extends DepositContent {
  const QrisContent(super.model, {super.key});

  @override
  State<QrisContent> createState() => _QrisContentState();
}

class _QrisContentState extends State<QrisContent> {
  final ssController = ScreenshotController();

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        const Text(
          "Screen capture atau screenshot QRIS dibawah ini dan scan di e-wallet atau m-banking Anda:",
        ),
        const SizedBox(height: 20),
        Center(
          child: Container(
              constraints: const BoxConstraints(
                maxWidth: 300,
              ),
              decoration: const BoxDecoration(
                borderRadius: BorderRadius.all(Radius.circular(8)),
                color: Colors.white,
              ),
              padding: const EdgeInsets.all(8),
              child: Column(
                children: [
                  Screenshot(
                      controller: ssController,
                      child: QrImageView(
                          backgroundColor: Colors.white,
                          data:
                              widget.model.paymentData["qrString"] ?? "NULL")),
                  const SizedBox(height: 4),
                  GestureDetector(
                    onTap: () async {
                      final image = await ssController.capture();
                      if (image != null) {
                        await ImageGallerySaver.saveImage(image);
                        if (mounted) {
                          Toast.success(context, "Kode QR disimpan ke galeri");
                        }
                      }
                    },
                    child: const Text(
                      "Simpan Kode QR",
                      style: TextStyle(
                        color: Colors.blue,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                  ),
                  const SizedBox(height: 8),
                ],
              )),
        ),
        const SizedBox(height: 20),
        Input(
          enabled: false,
          label: "Nominal Deposit",
          controller: TextEditingController(
              text: Meta.currencyFormatRp(widget.model.amount)),
        ),
        const SizedBox(height: 16),
        Input(
          enabled: false,
          label: "Biaya Admin",
          controller: TextEditingController(
              text: Meta.currencyFormatRp(widget.model.fee)),
        ),
        const SizedBox(height: 16),
        const Text(
          "Pastikan melakukan pembayaran sebelum kadaluarsa. Saldo yang diterima akan dipotong biaya admin.",
        ),
      ],
    );
  }
}
