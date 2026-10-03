// import http
import 'package:http/http.dart' as http;

class Internet {
  static Future<bool> isOnline() async {
    try {
      var response = await http.get(Uri.https(''));
      return response.statusCode == 200;
    } catch (e) {
      return false;
    }
  }

  static Future<http.Response> post(String url,
          {Map<String, String>? body, Map<String, String>? headers}) =>
      http.post(Uri.parse(url), body: body, headers: headers);

  static Future<http.Response> get(String url,
          {Map<String, String>? headers}) =>
      http.get(Uri.parse(url), headers: headers);
}
