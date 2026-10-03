import 'dart:convert';

import 'package:http/http.dart';

dynamic jsonDec(Response resp) {
  return jsonDecode(utf8.decode(resp.bodyBytes));
}
