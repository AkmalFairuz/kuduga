import 'package:finance_app/utils/meta.dart';
import 'package:finance_app/widgets/text/text_app_bar.dart';
import 'package:flutter/material.dart';
import 'package:webview_flutter/webview_flutter.dart';

class DepositWebViewContent extends StatefulWidget {
  const DepositWebViewContent({Key? key, required this.url}) : super(key: key);
  final String url;

  @override
  State<DepositWebViewContent> createState() => _DepositWebViewContentState();
}

class _DepositWebViewContentState extends State<DepositWebViewContent> {
  final WebViewController _controller = WebViewController();
  bool _loaded = false;
  bool _forcePop = false;

  @override
  void initState() {
    super.initState();

    _controller.setJavaScriptMode(JavaScriptMode.unrestricted);
    _controller.setNavigationDelegate(NavigationDelegate(onPageStarted: (url) {
      if (url.contains("example.invalid/return-payment")) {
        Navigator.pop(context);
      }
    }, onPageFinished: (url) {
      setState(() {
        _loaded = true;
      });
      Meta.webviewInjectFont(_controller);
      if (widget.url.contains("passport.duitku.com")) {
        _controller.runJavaScript(
            "document.getElementById('back').style.display = 'none';");
        _controller.runJavaScript(
            "document.getElementById('Image1').style.display = 'none';");
      }
    }));
    _controller.loadRequest(Uri.parse(widget.url));
  }

  @override
  Widget build(BuildContext context) {
    return WillPopScope(
        onWillPop: () async {
          if (_forcePop) {
            return true;
          }
          if (await _controller.canGoBack()) {
            _controller.goBack();
            return false;
          }
          return true;
        },
        child: Scaffold(
            appBar: AppBar(
              title: const AppBarTitle("Detail Pembayaran"),
              automaticallyImplyLeading: false,
              leading: IconButton(
                icon: const Icon(Icons.close),
                onPressed: () {
                  _forcePop = true;
                  Navigator.pop(context);
                },
              ),
            ),
            body: _loaded
                ? SizedBox(
                    width: double.infinity,
                    height: double.infinity,
                    child: WebViewWidget(controller: _controller))
                : const LinearProgressIndicator()));
  }
}
