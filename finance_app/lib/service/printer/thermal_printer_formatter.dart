import 'dart:math';

class ThermalPrinterFormatter {
  String _format = "";
  int _maxCharactersPerLine = 32;

  ThermalPrinterFormatter();

  String print() {
    return _format;
  }

  void setMaxCharactersPerLine(int max) {
    _maxCharactersPerLine = max;
  }

  void rawAppend(String v) {
    _format += v;
  }

  void _appendFormatted(String v,
      {required String prefix, required String suffix}) {
    final parts = <String>[];
    var tmp = "";
    for (final word in v.split(" ")) {
      if (word.length > _maxCharactersPerLine) {
        var idx = _maxCharactersPerLine - tmp.length;
        tmp += word.substring(0, idx);
        parts.add(tmp);
        tmp = "";
        while (idx < word.length) {
          tmp +=
              "${word.substring(idx, min(idx + _maxCharactersPerLine, word.length))} ";
          if (tmp.length < _maxCharactersPerLine) {
            break;
          }
          parts.add(tmp);
          tmp = "";
          idx += _maxCharactersPerLine;
        }
        continue;
      }
      if ("$tmp$word".length <= _maxCharactersPerLine) {
        tmp += "$word ";
        if (tmp.trim().length == _maxCharactersPerLine) {
          parts.add(tmp.trim());
          tmp = "";
        }
      } else {
        parts.add(tmp);
        tmp = "$word ";
      }
    }
    if (tmp.isNotEmpty) {
      parts.add(tmp);
    }
    int i = 0;
    for (final part in parts) {
      i++;
      rawAppend("$prefix${part.trim()}$suffix");
      if (i != parts.length) {
        line();
      }
    }
  }

  void appendFormatted(String v, {String prefix = "", String suffix = ""}) {
    var i = 0;
    final parts = v.split("\n");
    for (final line in parts) {
      i++;
      _appendFormatted(line, prefix: prefix, suffix: suffix);
      if (i != parts.length) {
        rawAppend("\n");
      }
    }
  }

  void stripes() {
    append(_strRepeat("=", _maxCharactersPerLine));
  }

  void append(String v) {
    // escape string
    var v2 = v
        .replaceAll("<", "?")
        .replaceAll(">", "?")
        .replaceAll("[", "?")
        .replaceAll("]", "?")
        .replaceAll("\n", " ");
    _format += v2;
  }

  void line({int len = 1}) {
    rawAppend(_strRepeat("\n", len));
  }

  String _strRepeat(String str, int len) {
    var ret = "";
    for (int i = 0; i < len; i++) {
      ret += str;
    }
    return ret;
  }
}
