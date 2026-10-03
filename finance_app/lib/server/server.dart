import 'dart:convert';
import 'dart:io';

import 'package:device_info_plus/device_info_plus.dart';
import 'package:finance_app/constants.dart';
import 'package:finance_app/controller/account_controller.dart';
import 'package:finance_app/screens/intro_screen.dart';
import 'package:finance_app/service/auth_service.dart';
import 'package:finance_app/utils/alert.dart';
import 'package:finance_app/utils/screens.dart';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'package:web_socket_channel/web_socket_channel.dart';

import '../utils/local_storage.dart';

class Server {
  static String? _token;
  static bool? _isAuthenticatedCache;
  static bool loaded = false;
  static String _userAgent = "${Constants.appName}/${Constants.version};";
  static bool _initiated = false;

  static void init() async {
    if (_initiated) {
      return;
    }
    _initiated = true;
    var deviceInfo = DeviceInfoPlugin();
    if (Platform.isAndroid) {
      var info = await deviceInfo.androidInfo;
      _userAgent += " Android/${info.version.release} (${info.model};)";
    } else if (Platform.isIOS) {
      var info = await deviceInfo.iosInfo;
      _userAgent += " iOS/${info.systemVersion} (${info.model})";
    }
  }

  static Future<bool> setToken(String token) async {
    _token = token;

    await LocalStorage.secureSet('authorizationToken', token);

    return true;
  }

  static Future<void> logout(BuildContext? context) async {
    if (context != null) {
      Alert.showLoading(context: context, message: 'Logging out...');
    }

    try {
      if (_token != null) {
        await AuthService.logout();
      }
    } catch (ignored) {
      // ignore error
    }

    _token = null;
    _isAuthenticatedCache = null;
    await LocalStorage.secureSet('authorizationToken', '');
    AccountController.destroyInstance();

    if (context != null && context.mounted) {
      Alert.closeLoading(context: context);
    }
    Screens.replaceAll(const IntroScreen());
  }

  static bool hasToken() {
    return _token != null;
  }

  static Future<void> loadToken() async {
    const authField = 'authorizationToken';
    _token = await LocalStorage.secureGet(authField);
    if (_token == null) {
      final tmpTok = await LocalStorage.get(authField);
      if (tmpTok != null) {
        await LocalStorage.set(authField, '');
        await LocalStorage.secureSet(authField, tmpTok);
        _token = tmpTok;
      }
    }
    if (_token == "") {
      _token = null;
    }
  }

  static Future<http.Response> get(String path,
          {Map<String, String>? headers, Map<String, dynamic>? query}) =>
      http
          .get(_parse(path, query), headers: _getHeaders(headers))
          .then(_processResponse);

  static Future<http.Response> post(
    String path, {
    Object? body,
    Map<String, String>? headers,
    Map<String, dynamic>? query,
    List<http.MultipartFile> files = const [],
  }) {
    if (files.isNotEmpty) {
      var req = http.MultipartRequest("POST", _parse(path, query));
      req.headers.addAll(_getHeaders(headers));
      req.fields.addAll(body as Map<String, String>);
      req.files.addAll(files);
      return req.send().then((res) async {
        return await http.Response.fromStream(res);
      }).then(_processResponse);
    }
    return http
        .post(_parse(path, query), body: body, headers: _getHeaders(headers))
        .then(_processResponse);
  }

  static Future<http.Response> put(String path,
          {Object? body,
          Map<String, String>? headers,
          Map<String, dynamic>? query}) =>
      http
          .put(_parse(path, query), body: body, headers: _getHeaders(headers))
          .then(_processResponse);

  static Future<http.Response> delete(String path,
          {Map<String, String>? headers, Map<String, dynamic>? query}) =>
      http
          .delete(_parse(path, query), headers: _getHeaders(headers))
          .then(_processResponse);

  static Future<http.Response> patch(String path,
          {Object? body,
          Map<String, String>? headers,
          Map<String, dynamic>? query}) =>
      http
          .patch(_parse(path, query), body: body, headers: _getHeaders(headers))
          .then(_processResponse);

  static Future<http.Response> _processResponse(http.Response response) async {
    if (loaded && response.statusCode == 401) {
      if (response.headers.containsKey('x-require-auth') &&
          response.headers['x-require-auth'] == 'true') {
        Alert.message('Sesi tidak valid. Silakan login kembali.')
            .then((_) => logout(null));
        throw Exception('Unauthorized');
      }
    }
    if (response.statusCode >= 400) {
      throw ServerRequestError.init(response);
    }
    return response;
  }

  static Map<String, String> _getHeaders(Map<String, String>? headers) {
    headers ??= {};
    headers["user-agent"] = _userAgent;
    headers.addAll(_getAuthHeaders());
    return headers;
  }

  static Map<String, String> _getAuthHeaders() {
    Map<String, String> headers = {};
    if (_token != null) {
      headers['Authorization'] = 'Bearer $_token';
    }
    return headers;
  }

  static Future<bool> isAuthenticated(
      {bool cached = false, bool init = false}) async {
    if (cached && _isAuthenticatedCache != null) {
      return _isAuthenticatedCache!;
    }
    if (_token == null) {
      return false;
    }
    try {
      await get('/auth/details');
      return _isAuthenticatedCache = true;
    } on ServerRequestError catch (e) {
      if (e.response.statusCode < 500) {
        return _isAuthenticatedCache = false;
      }
      rethrow;
    }
  }

  static Uri _parse(String path, Map<String, dynamic>? query) {
    // convert all query values to string
    query = query?.map((key, value) => MapEntry(key, value.toString()));
    if (Constants.secureApi) {
      return Uri.https(Constants.apiBaseUrl, path, query);
    }
    return Uri.http(Constants.apiBaseUrl, path, query);
  }

  static WebSocketChannel createWebSocketChannel(String uri) {
    final channel = WebSocketChannel.connect(Uri.parse(
        "${Constants.secureApi ? "wss" : "ws"}://${Constants.apiBaseUrl}$uri"));
    channel.sink.add(jsonEncode({"authorization": "Bearer ${_token ?? ""}"}));
    return channel;
  }
}

class ServerRequestError implements Exception {
  final String msg;
  final http.Response response;
  final int errorCode;

  const ServerRequestError(this.msg, this.response, this.errorCode);

  factory ServerRequestError.init(http.Response response) {
    int errorCode = int.parse(response.headers['x-finance-error-code'] ?? '0');
    if (response.headers.containsKey('content-type') &&
        response.headers['content-type'] == 'application/json') {
      try {
        var decoded = jsonDecode(response.body) as Map<String, dynamic>;
        if (decoded.containsKey('message')) {
          return ServerRequestError(decoded['message'], response, errorCode);
        }
      } on FormatException catch (e) {
        // ignored
      }
    }

    return ServerRequestError(
        response.reasonPhrase ?? 'Error code: ${response.statusCode}',
        response,
        errorCode);
  }

  @override
  String toString() => msg;
}
