import 'dart:convert';
import 'package:web_socket_channel/io.dart';
import 'package:scoped_model/scoped_model.dart';

class MyModel extends Model {
  bool _isConnected = false;
  IOWebSocketChannel? _channel;

  bool get isConnected => _isConnected;

  MyModel() {
    connectWebSock();
  }

  void onConnect() {
    _isConnected = true;
    notifyListeners();
  }

  void onDisconnect() {
    _isConnected = false;
    notifyListeners();
  }

  Future<void> connectWebSock() async {
    _channel?.sink.close();
    _channel = IOWebSocketChannel.connect(
      Uri.parse('wss://ts.beebs.dev/ws'),
      pingInterval: const Duration(seconds: 10),
      headers: {},
    );
    onConnect();
    _channel!.stream.listen(
      (message) => handleWebSockMessage(message),
      onError: (err) async {
        print('websock error: $err');
        onDisconnect();
        connectWebSock();
      },
      onDone: () async {
        onDisconnect();
        connectWebSock();
      },
    );
  }

  void handleWebSockMessage(dynamic message) {
    final obj = jsonDecode(message) as Map<String, dynamic>;
    switch (obj['type']) {
      case 'ping':
        _channel?.sink.add(jsonEncode({'type': 'pong'}));
        break;
      default:
        break;
    }
  }
}
