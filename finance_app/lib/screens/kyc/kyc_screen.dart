import 'package:camera/camera.dart';
import 'package:finance_app/screens/utility/image_viewer_screen.dart';
import 'package:finance_app/service/user_service.dart';
import 'package:finance_app/styles/color_helper.dart';
import 'package:finance_app/utils/alert.dart';
import 'package:finance_app/utils/meta.dart';
import 'package:finance_app/utils/screens.dart';
import 'package:finance_app/widgets/button/button.dart';
import 'package:finance_app/widgets/input/input.dart';
import 'package:finance_app/widgets/text/text_app_bar.dart';
import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:image/image.dart' as img;

class KycScreen extends StatelessWidget {
  const KycScreen({Key? key}) : super(key: key);

  Widget _buildRequired(IconData icon, String title, String subtitle) {
    return Row(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Icon(
          icon,
          size: 32,
        ),
        const SizedBox(width: 16),
        Expanded(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                title,
                style:
                    const TextStyle(fontSize: 16, fontWeight: FontWeight.w600),
              ),
              const SizedBox(height: 4),
              Text(
                subtitle,
                style: TextStyle(color: Colors.grey[600]),
              ),
            ],
          ),
        ),
      ],
    );
  }

  void _startKyc() async {
    var result = await Screens.to<Uint8List>(const _KycStep1Screen());
    if (result == null) {
      return;
    }
    Screens.to(_KycStep2Screen(documentImage: result));
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: const TextAppBar("Verifikasi Identitas"),
      body: Container(
        padding: const EdgeInsets.symmetric(horizontal: 16),
        width: double.infinity,
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            const Text("Verifikasi Identitas",
                style: TextStyle(fontSize: 24, fontWeight: FontWeight.bold)),
            const SizedBox(height: 16),
            Text(
                "Untuk memverifikasi identitas anda, mohon untuk mengirimkan dokumen yang diperlukan.",
                textAlign: TextAlign.center,
                style: TextStyle(color: Colors.grey[600])),
            const SizedBox(height: 32),
            _buildRequired(
                Icons.person_pin_rounded,
                "Mengambil foto kartu identitas anda",
                "Mohon untuk menyiapkan KTP"),
            const SizedBox(height: 16),
            _buildRequired(
                Icons.edit_note,
                "Mengisi formulir tentang identitas anda",
                "Mengisi nama lengkap dan NIK yang sesuai dengan KTP"),
            const SizedBox(height: 40),
            Button(onPressed: _startKyc, child: const Text("Mulai")),
          ],
        ),
      ),
    );
  }
}

Uint8List? _cropImage(Uint8List imageBytes) {
  var originalImage = img.decodeImage(imageBytes);
  if (originalImage == null) {
    return null;
  }
  const double aspectRatio = 85.6 / 54;
  final int originalWidth = originalImage.width;
  final int originalHeight = originalImage.height;

  int newWidth, newHeight, x, y;

  if (originalWidth / originalHeight > aspectRatio) {
    newHeight = originalHeight;
    newWidth = (newHeight * aspectRatio).toInt();
    x = (originalWidth - newWidth) ~/ 2;
    y = 0;
  } else {
    newWidth = originalWidth;
    newHeight = newWidth ~/ aspectRatio;
    x = 0;
    y = (originalHeight - newHeight) ~/ 2;
  }

  return img.encodePng(img.copyCrop(originalImage,
      x: x, y: y, width: newWidth, height: newHeight));
}

class _KycStep1Screen extends StatefulWidget {
  const _KycStep1Screen({Key? key}) : super(key: key);

  @override
  State<_KycStep1Screen> createState() => _KycStep1State();
}

class _KycStep1State extends State<_KycStep1Screen>
    with WidgetsBindingObserver {
  late CameraController cameraController;
  late List<CameraDescription> cameras;
  bool cameraOn = false;
  bool hasShot = false;
  Uint8List? pictureResult;

  @override
  void initState() {
    super.initState();

    loadCamera();
  }

  @override
  void dispose() {
    cameraController.dispose();

    super.dispose();
  }

  void loadCamera() async {
    pictureResult = null;
    cameras = await availableCameras();
    if (cameras.isEmpty) {
      await Alert.message("Kamera tidak tersedia");
      Screens.back();
      return;
    }
    cameraController = CameraController(
        cameras
                .where((cam) => cam.lensDirection == CameraLensDirection.back)
                .firstOrNull ??
            cameras.first,
        ResolutionPreset.veryHigh,
        enableAudio: false);
    await cameraController.initialize();
    cameraController.setFlashMode(FlashMode.off);
    setState(() {
      cameraOn = true;
    });
  }

  @override
  void didChangeAppLifecycleState(AppLifecycleState state) {
    if (!cameraController.value.isInitialized) {
      return;
    }

    if (state == AppLifecycleState.inactive) {
      cameraController.dispose();
    } else if (state == AppLifecycleState.resumed) {
      loadCamera();
    }
  }

  void _onShot() async {
    Alert.withLoading((done) async {
      var result = await cameraController.takePicture();
      var imageBytes = await result.readAsBytes();
      pictureResult = await compute(_cropImage, imageBytes);
      if (pictureResult == null) {
        done();
        return;
      }
      cameraOn = false;
      cameraController.dispose();
      setState(() {});
      done();
    }, message: "Mengambil foto...");
  }

  Widget _buildShot() {
    return GestureDetector(
      onTap: _onShot,
      child: Container(
        width: 80.sp,
        height: 80.sp,
        decoration: BoxDecoration(
            borderRadius: BorderRadius.circular(999), color: Meta.color),
        child: Center(
          child: Icon(
            Icons.camera_alt,
            size: 32.sp,
            color: Colors.white,
          ),
        ),
      ),
    );
  }

  // TODO: rewrite this, very hacky code
  Widget _buildContent() {
    if (!cameraOn) {
      if (pictureResult != null) {
        return Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Image.memory(pictureResult!),
            const SizedBox(height: 50),
            Row(
              children: [
                Expanded(
                  child: Button(
                      onPressed: () {
                        loadCamera();
                      },
                      variant: ButtonVariant.outline,
                      child: const Text("Ulangi")),
                ),
                const SizedBox(width: 16),
                Expanded(
                  child: Button(
                      onPressed: () {
                        Navigator.of(context).pop(pictureResult);
                      },
                      child: const Text("Lanjut")),
                )
              ],
            )
          ],
        );
      }
      return const Center(
        child: CircularProgressIndicator(),
      );
    }
    return Stack(
      alignment: Alignment.center,
      children: [
        CameraPreview(cameraController),
        ColorFiltered(
          colorFilter:
              ColorFilter.mode(ColorHelper.bg(context), BlendMode.srcOut),
          child: Stack(alignment: Alignment.center, children: [
            Container(
              decoration: const BoxDecoration(
                  color: Colors.black,
                  backgroundBlendMode: BlendMode
                      .dstOut), // This one will handle background + difference out
            ),
            Column(
              mainAxisAlignment: MainAxisAlignment.end,
              children: [
                Expanded(
                    child: Column(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    AspectRatio(
                      aspectRatio: 85.6 / 54,
                      child: Container(
                        color: Colors.red,
                      ),
                    ),
                  ],
                )),
              ],
            ),
          ]),
        ),
        Column(
          children: [
            SizedBox(height: 32.sp),
            const Text(
              "Posisikan KTP pada kotak dibawah. Pastikan hasil foto terlihat jelas.",
              textAlign: TextAlign.center,
            )
          ],
        ),
        Column(
          mainAxisAlignment: MainAxisAlignment.end,
          children: [
            _buildShot(),
            SizedBox(height: 40.sp),
          ],
        )
      ],
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: const TextAppBar("Foto Kartu Identitas"),
      body: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Expanded(
            child: Padding(
              padding: const EdgeInsets.symmetric(horizontal: 16),
              child: _buildContent(),
            ),
          )
        ],
      ),
    );
  }
}

class _KycStep2Screen extends StatefulWidget {
  const _KycStep2Screen({Key? key, required this.documentImage})
      : super(key: key);

  final Uint8List documentImage;

  @override
  State<_KycStep2Screen> createState() => _KycStep2State();
}

class _KycStep2State extends State<_KycStep2Screen> {
  final nikController = TextEditingController();
  final fullNameController = TextEditingController();

  void _onContinue() async {
    Alert.withLoading((done) async {
      var result = await UserService.submitKyc(
          documentId: nikController.text,
          fullName: fullNameController.text,
          ktpData: widget.documentImage);
      done();
      await Alert.message(result);
      if (mounted) {
        Navigator.of(context).pop(true);
      }
    });
  }

  @override
  Widget build(BuildContext context) {
    final img = MemoryImage(widget.documentImage);
    return Scaffold(
      appBar: const TextAppBar("Informasi Pribadi"),
      body: SingleChildScrollView(
          child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text("Foto KTP", style: Input.labelTextStyle(context)),
            const SizedBox(height: 4),
            GestureDetector(
              onTap: () => Screens.to(ImageViewerScreen(provider: img)),
              child: Image(image: img, width: double.infinity),
            ),
            const SizedBox(height: 16),
            Input(
              label: "Nama Lengkap",
              controller: fullNameController,
            ),
            const SizedBox(height: 16),
            Input(label: "NIK (16 digit)", controller: nikController),
            const SizedBox(height: 16),
            const Text("Verifikasi membutuhkan waktu selama 1 hari kerja."),
            const SizedBox(height: 16),
            Button(
                onPressed: _onContinue,
                width: double.infinity,
                child: const Text("Kirim"))
          ],
        ),
      )),
    );
  }
}
