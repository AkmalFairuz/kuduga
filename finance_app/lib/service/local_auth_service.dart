import 'package:finance_app/utils/local_storage.dart';
import 'package:local_auth/local_auth.dart';

class LocalAuthService {
  static final _fingerprintLocker = LocalAuthentication();

  static Future<bool> isFingerprintLockEnabled() async {
    return await LocalStorage.get("fingerprintLockEnable") ?? false;
  }

  static Future<bool> isFingerprintAvailable() async {
    final List<BiometricType> availableBiometrics =
        await _fingerprintLocker.getAvailableBiometrics();
    return availableBiometrics.isNotEmpty;
  }

  static Future<void> setFingerprintLockEnabled(bool enabled) async {
    if (await isFingerprintLockEnabled() == enabled) {
      return;
    }
    final List<BiometricType> availableBiometrics =
        await _fingerprintLocker.getAvailableBiometrics();

    if (availableBiometrics.isEmpty) {
      throw Exception(
          "Autentikasi menggunakan sidik jari atau muka tidak tersedia di perangkat Anda");
    }

    if (!(await tryFingerprintAuth(force: true))) {
      throw Exception("Autentikasi gagal");
    }
    await LocalStorage.set("fingerprintLockEnable", enabled);
  }

  static Future<bool> tryFingerprintAuth({bool force = false}) async {
    if ((!(await isFingerprintLockEnabled()) ||
            !(await isFingerprintAvailable())) &&
        !force) {
      return true;
    }
    return await _fingerprintLocker.authenticate(
        localizedReason: "Mohon lakukan autentikasi untuk melanjutkan",
        options: const AuthenticationOptions(biometricOnly: true));
  }
}
