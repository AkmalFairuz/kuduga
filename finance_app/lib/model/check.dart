class CheckResultModel {
  final String name;
  final String value;

  CheckResultModel(this.name, this.value);

  factory CheckResultModel.fromJson(List<dynamic> json) {
    return CheckResultModel(json[0] as String, json[1] as String);
  }
}
