import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:shared_preferences/shared_preferences.dart';

class LocalStorage {
  static SharedPreferences? _sharedPreferences;

  static FlutterSecureStorage? _secureStorage;

  static Future<SharedPreferences> _getSharedPreferences() async {
    _sharedPreferences ??= await SharedPreferences.getInstance();
    return _sharedPreferences!;
  }

  static Future<FlutterSecureStorage> _getSecureStorage() async {
    _secureStorage ??= const FlutterSecureStorage(
      aOptions: AndroidOptions(
        encryptedSharedPreferences: true,
      ),
    );
    return _secureStorage!;
  }

  static Future<void> secureSet(String key, dynamic value) async {
    final s = await _getSecureStorage();
    return s.write(key: key, value: value);
  }

  static Future<String?> secureGet(String key) async {
    final s = await _getSecureStorage();
    return s.read(key: key);
  }

  static Future<bool> set(String key, dynamic value) async {
    SharedPreferences sharedPreferences = await _getSharedPreferences();
    if (value is String) {
      return sharedPreferences.setString(key, value);
    } else if (value is int) {
      return sharedPreferences.setInt(key, value);
    } else if (value is double) {
      return sharedPreferences.setDouble(key, value);
    } else if (value is bool) {
      return sharedPreferences.setBool(key, value);
    } else if (value is List<String>) {
      return sharedPreferences.setStringList(key, value);
    } else {
      return Future.value(false);
    }
  }

  static Future<T?> get<T>(String key, {T? def}) async {
    SharedPreferences sharedPreferences = await _getSharedPreferences();
    T? value = sharedPreferences.get(key) as T?;
    if (value == null) {
      return def;
    }
    return value;
  }

  static Future<bool> isExists(String key) async {
    SharedPreferences sharedPreferences = await _getSharedPreferences();
    return sharedPreferences.containsKey(key);
  }

  static Future<bool> remove(String key) async {
    SharedPreferences sharedPreferences = await _getSharedPreferences();
    return sharedPreferences.remove(key);
  }
}
