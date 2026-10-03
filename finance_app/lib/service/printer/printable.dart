class ThermalPrintable {
  String thermalPrint(int maxCharacterPerLine) {
    return "";
  }
}

class ThermalPrintableText implements ThermalPrintable {
  const ThermalPrintableText(this.text);

  final String text;

  @override
  String thermalPrint(int maxCharacterPerLine) {
    return text;
  }
}
