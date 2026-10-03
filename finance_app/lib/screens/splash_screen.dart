import 'package:finance_app/constants.dart';
import 'package:finance_app/controller/account_controller.dart';
import 'package:finance_app/screens/home_screen.dart';
import 'package:finance_app/screens/intro_screen.dart';
import 'package:finance_app/server/server.dart';
import 'package:finance_app/service/app_service.dart';
import 'package:finance_app/service/local_auth_service.dart';
import 'package:finance_app/service/notification_service.dart';
import 'package:finance_app/utils/alert.dart';
import 'package:finance_app/utils/meta.dart';
import 'package:finance_app/utils/screens.dart';
import 'package:finance_app/widgets/icon/brand_icon.dart';
import 'package:firebase_core/firebase_core.dart';
import 'package:flutter/material.dart';
import 'package:flutter_udid/flutter_udid.dart';
import 'package:url_launcher/url_launcher.dart';

class SplashScreen extends StatelessWidget {
  const SplashScreen({Key? key, this.onComplete}) : super(key: key);

  final VoidCallback? onComplete;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
        body: Container(
      color: Meta.color,
      padding: const EdgeInsets.all(20),
      width: double.infinity,
      child: const _SplashScreen(),
    ));
  }
}

class _SplashScreen extends StatefulWidget {
  const _SplashScreen({Key? key}) : super(key: key);

  @override
  _SplashScreenState createState() => _SplashScreenState();
}

class _SplashScreenState extends State<_SplashScreen> {
  Widget _targetScreen = const IntroScreen();
  String _message = "Initializing...";

  @override
  void initState() {
    super.initState();

    _startLoading();
  }

  Future<void> _initializeStuff() async {
    await Firebase.initializeApp();

    FlutterUdid.udid.then((value) => {Meta.udid = value});

    Server.init();
    NotificationService.instance.init();
  }

  void _startLoading() async {
    _setMessage('Loading...');

    while (true) {
      try {
        final ver = await AppService.getVersionInfo();
        if (Constants.apiVersion < ver.minimumVersion) {
          while (true) {
            if ((await Alert.message(
                    "Update diperlukan untuk memulai aplikasi.\nMohon untuk update aplikasi Kuduga di Play Store.",
                    btnMsg: "Update Sekarang")) ??
                false) {
              await launchUrl(
                  Uri.parse("market://details?id=${Constants.packageName}"),
                  mode: LaunchMode.externalApplication);
            }
            await Future.delayed(const Duration(seconds: 2));
          }
        }
        break;
      } catch (e) {
        await Alert.message(e.toString());
      }
    }

    Future.delayed(const Duration(milliseconds: 500), () {
      _onLoading();
    });
  }

  void _onLoading() async {
    await _initializeStuff();
    try {
      _setMessage('Loading token...');
      await Server.loadToken();
      _setMessage('Checking authentication...');
      bool isAuthenticated = await Server.isAuthenticated();
      if (isAuthenticated) {
        _targetScreen = const HomeScreen();
      }
      _setMessage('Launching app...');
      _onComplete();
    } catch (e) {
      String msg = e.toString();
      if (msg.contains('Connection refused')) {
        msg = 'Koneksi ke server gagal.';
      }
      Alert.show(title: 'Terjadi kesalahan', message: msg, actions: [
        TextButton(
          onPressed: () {
            _startLoading();
            Screens.back();
          },
          child: const Text("Coba lagi"),
        ),
      ]);
    }
  }

  void _setMessage(String message) {
    setState(() {
      _message = message;
    });
  }

  void _onComplete() async {
    Server.loaded = true;
    if (Server.hasToken()) {
      await AccountController.getInstance().fetchAll();

      _setMessage("Authenticating...");
      while (!await LocalAuthService.tryFingerprintAuth()) {
        await Alert.message("Autentikasi gagal, jika Anda memiliki masalah "
            "dengan autentikasi, mohon nonaktifkan kunci fingerprint di "
            "pengaturan sistem");
      }
    }
    Screens.replaceAll(_targetScreen);
  }

  @override
  Widget build(BuildContext context) {
    return Column(
      mainAxisAlignment: MainAxisAlignment.center,
      crossAxisAlignment: CrossAxisAlignment.center,
      children: [
        const Expanded(
            child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            BrandIcon(size: 280),
            SizedBox(height: 48),
            Text(
              "kuduga",
              style: TextStyle(
                  color: Colors.white,
                  fontSize: 50,
                  fontWeight: FontWeight.bold),
            )
          ],
        )),
        // Loading indicator
        SizedBox(
          height: 80,
          width: double.infinity,
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            crossAxisAlignment: CrossAxisAlignment.center,
            children: [
              SizedBox(
                width: 20,
                height: 20,
                child: CircularProgressIndicator(
                  strokeWidth: 2,
                  color: Colors.grey[50],
                ),
              ),
              const SizedBox(height: 12),
              Text(
                _message,
                style: TextStyle(
                  fontSize: 12,
                  color: Colors.grey[50],
                ),
              ),
            ],
          ),
        )
      ],
    );
  }
}
