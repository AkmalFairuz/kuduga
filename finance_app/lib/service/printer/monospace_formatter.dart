import 'dart:math';

class MonospaceFormatter {
  final int nameLength;
  final int valueLength;
  final String separator;
  final List<_MonospaceItemData> _data = [];

  MonospaceFormatter({
    this.nameLength = 16,
    this.valueLength = 16,
    this.separator = " : ",
  });

  void add(String name, String value) {
    _data.add(_MonospaceItemData(name, value));
  }

  String format() {
    var nameLength2 = 0;
    for (final data in _data) {
      nameLength2 = max(nameLength2, data.name.length);
    }
    nameLength2 = min(nameLength2, nameLength);

    String ret = "";
    var sepEmpty = _repeat(" ", separator.length);
    var nameEmpty = _repeat(" ", nameLength2);
    var valueEmpty = _repeat(" ", valueLength);

    for (final data in _data) {
      final nameParts = _formatChild(data.name, nameLength2);
      final valueParts = _formatChild(data.value, valueLength);
      for (int i = 0; i < max(nameParts.length, valueParts.length); i++) {
        final namePart = nameParts.length > i ? nameParts[i] : nameEmpty;
        final valuePart = valueParts.length > i ? valueParts[i] : valueEmpty;
        ret += "$namePart${i == 0 ? separator : sepEmpty}$valuePart\n";
      }
    }

    return ret.substring(0, ret.length - 1);
  }

  List<String> _formatChild(String text, int maxLen) {
    text = text.replaceAll("  ", " ");
    if (text.length <= maxLen) {
      return [_pad(text, maxLen)];
    }
    final ret = <String>[];
    final words = text.split(" ");
    var tmp = "";
    for (final word in words) {
      if (word.length > maxLen) {
        var idx = maxLen - tmp.length;
        tmp += word.substring(0, idx);
        ret.add(tmp);
        tmp = "";
        while (idx < word.length) {
          tmp += "${word.substring(idx, min(idx + maxLen, word.length))} ";
          if (tmp.length < maxLen) {
            break;
          }
          ret.add(tmp);
          tmp = "";
          idx += maxLen;
        }
        continue;
      }
      if ("$tmp$word".length <= maxLen) {
        tmp += "$word ";
        if (tmp.trim().length == maxLen) {
          ret.add(tmp.trim());
          tmp = "";
        }
      } else {
        ret.add(tmp);
        tmp = "$word ";
      }
    }
    if (tmp.isNotEmpty) {
      ret.add(tmp);
    }
    return ret.map((e) => _pad(e.trim(), maxLen)).toList();
  }

  String _pad(String text, int len) {
    if (text.length >= len) {
      return text;
    }
    return "$text${_repeat(" ", len - text.length)}";
  }

  String _repeat(String str, int len) {
    var ret = "";
    for (int i = 0; i < len; i++) {
      ret += str;
    }
    return ret;
  }
}

class _MonospaceItemData {
  final String name;
  final String value;

  const _MonospaceItemData(this.name, this.value);
}
