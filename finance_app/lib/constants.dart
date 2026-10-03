class Constants {
  static const String appName = "Kuduga";
  static const String version = "1.2.0";
  static const int apiVersion = 7;
  static const String packageName = "com.kuduga.app";
  // Host and optional port only, without a scheme or path.
  static const String apiBaseUrl = String.fromEnvironment(
    'API_HOST',
    defaultValue: 'localhost:3001',
  );
  static const bool secureApi = bool.fromEnvironment(
    'API_SECURE',
    defaultValue: false,
  );
}
