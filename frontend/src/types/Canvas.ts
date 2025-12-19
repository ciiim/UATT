export interface CanvasComponent {
  ID: string;
  Type: string;
  Label: string;
  Position: { X: number; Y: number };
  Value?: string; // 文字框的输出值
  AttachApp?: string; // 按钮的操作类型
}

export interface Connection {
  FromID: string; // 按钮ID
  ToID: string;   // 文字框ID
}