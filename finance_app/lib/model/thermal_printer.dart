class ThermalPrinterConnection {
  String name;
  String address;

  ThermalPrinterConnection(this.name, this.address);

  @override
  String toString() {
    return "$name@$address";
  }
}
