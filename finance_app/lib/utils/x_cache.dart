import 'dart:convert';
import 'dart:typed_data';

import 'package:crypto/crypto.dart';
import 'package:flutter_cache_manager/flutter_cache_manager.dart';

class XCache {
  static Future<String> network(String url) async {
    final file = await DefaultCacheManager().getSingleFile(url);
    final bytes = await file.readAsBytes();
    return String.fromCharCodes(bytes);
  }

  static String _hashKey(String key) {
    return "x_${sha1.convert(utf8.encode(key))}";
  }

  static Future<void> set(String key, String val) async {
    await DefaultCacheManager().putFile(
      _hashKey(key),
      Uint8List.fromList(val.codeUnits),
    );
  }

  static Future<String?> get(String key) async {
    final file = await DefaultCacheManager().getFileFromCache(_hashKey(key));
    if (file == null) {
      return null;
    }
    final bytes = await file.file.readAsBytes();
    return String.fromCharCodes(bytes);
  }
}
