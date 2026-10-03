import 'package:finance_app/utils/meta.dart';
import 'package:finance_app/widgets/text/text_app_bar.dart';
import 'package:flutter/material.dart';
import 'package:webview_flutter/webview_flutter.dart';

enum WebviewScreenType {
  url,
  html,
}

class WebviewScreen extends StatefulWidget {
  const WebviewScreen(this.data, {Key? key, this.type = WebviewScreenType.url})
      : super(key: key);

  final String data;
  final WebviewScreenType type;

  @override
  State<WebviewScreen> createState() => _WebviewScreenState();
}

class _WebviewScreenState extends State<WebviewScreen> {
  final WebViewController controller = WebViewController();
  String title = "";
  bool loaded = false;

  @override
  void initState() {
    super.initState();

    controller.setJavaScriptMode(JavaScriptMode.unrestricted);
    controller.setBackgroundColor(Colors.grey[50]!);
    controller.setNavigationDelegate(NavigationDelegate(
      onPageFinished: (_) async {
        Meta.webviewInjectFont(controller);
        title = (await controller.getTitle()) ?? "";
        loaded = true;
        setState(() {});
      },
    ));
    if (widget.type == WebviewScreenType.html) {
      controller.loadHtmlString(widget.data);
      return;
    }
    controller.loadRequest(Uri.parse(widget.data));
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: TextAppBar(title),
      body: SizedBox(
        height: double.infinity,
        width: double.infinity,
        child: loaded
            ? WebViewWidget(
                controller: controller,
              )
            : const Center(
                child: CircularProgressIndicator(),
              ),
      ),
    );
  }
}
