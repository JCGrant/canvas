import * as React from 'react';

const canvasStyle = {
  border: '1px solid black',
};

interface IPaintMessage {
  last: { x: number; y: number };
  current: { x: number; y: number };
  brushSize: number;
  brushColor: string;
}

interface IPaintCanvasState {
  brushColor: string;
  brushSize: number;
  lastPoint: {
    x: number;
    y: number;
  };
  mouseDown: boolean;
}

const stateKey = 'paint-state';

const defaultState: IPaintCanvasState = {
  brushColor: '#000000',
  brushSize: 10,
  lastPoint: { x: 0, y: 0 },
  mouseDown: false,
};

class PaintCanvas extends React.Component<{}, IPaintCanvasState> {
  private canvas: HTMLCanvasElement;
  private ctx: CanvasRenderingContext2D;

  private socket: WebSocket;

  constructor(props: {}) {
    super(props);
    const storedStateStr = localStorage.getItem(stateKey);
    const storedState = storedStateStr ? JSON.parse(storedStateStr) : {};
    this.state = {
      ...defaultState,
      ...storedState,
    };
    this.socket = this.setupWebsocket();
  }

  public render() {
    return (
      <div>
        <div>
          <input
            type="color"
            value={this.state.brushColor}
            onChange={this.onChangeBrushColor}
          />
          <input
            type="range"
            min={1}
            max={100}
            value={this.state.brushSize}
            onChange={this.onChangeBrushSize}
          />
          <button onClick={this.onClickSave}>Save</button>
        </div>
        <canvas
          ref={this.canvasRef}
          style={canvasStyle}
          height={1000}
          width={1000}
          onMouseDown={this.onMouseDown}
          onMouseMove={this.onMouseMove}
        />
      </div>
    );
  }

  public componentDidMount() {
    window.addEventListener('mouseup', this.onMouseUp);
  }

  public componentDidUpdate() {
    localStorage.setItem(stateKey, JSON.stringify(this.state));
  }

  public componentWillUnmount() {
    window.removeEventListener('mouseup', this.onMouseUp);
  }

  private canvasRef = (canvas: HTMLCanvasElement) => {
    this.canvas = canvas;
    this.ctx = canvas.getContext('2d')!;
  };

  private onMouseMove = (e: React.MouseEvent) => {
    const currentPoint = {
      x: e.pageX - this.canvas.offsetLeft,
      y: e.pageY - this.canvas.offsetTop,
    };
    this.setState({ lastPoint: currentPoint });
    if (!this.state.mouseDown) {
      return;
    }
    this.paint(
      this.state.lastPoint,
      currentPoint,
      this.state.brushSize,
      this.state.brushColor
    );
    const data: IPaintMessage = {
      brushColor: this.state.brushColor,
      brushSize: this.state.brushSize,
      current: currentPoint,
      last: this.state.lastPoint,
    };
    try {
      this.socket.send(JSON.stringify(data));
    } catch (e) {
      // tslint:disable-next-line no-console
      console.error(e);
    }
  };

  private onMouseDown = (e: React.MouseEvent) => {
    e.persist();
    this.setState({ mouseDown: true }, () => this.onMouseMove(e));
  };

  private onMouseUp = () => {
    this.setState({ mouseDown: false });
  };

  private onChangeBrushColor = (e: React.ChangeEvent<HTMLInputElement>) => {
    const brushColor = e.target.value;
    this.setState({ brushColor });
  };

  private onChangeBrushSize = (e: React.ChangeEvent<HTMLInputElement>) => {
    const brushSize = parseInt(e.target.value, 10);
    this.setState({ brushSize });
  };

  private onClickSave = () => {
    const image = this.canvas
      .toDataURL('image/png')
      .replace('image/png', 'image/octet-stream');
    const link = document.createElement('a');
    link.setAttribute('download', 'painting.png');
    link.setAttribute('href', image);
    link.click();
  };

  private setupWebsocket = (): WebSocket => {
    const webSocketProtocol =
      window.location.protocol === 'http:' ? 'ws:' : 'wss:';
    const socket = new WebSocket(
      `${webSocketProtocol}//${window.location.host}/ws`
    );
    // tslint:disable-next-line no-console
    socket.onopen = () => console.log('connected');
    socket.onmessage = this.onMessage;
    return socket;
  };

  private onMessage = (e: { data: string }) => {
    const data: IPaintMessage = JSON.parse(e.data);
    this.paint(data.last, data.current, data.brushSize, data.brushColor);
  };

  private paint = (
    last: { x: number; y: number },
    current: { x: number; y: number },
    brushSize: number = 1,
    color: string = 'black'
  ) => {
    this.ctx.strokeStyle = color;
    this.ctx.lineWidth = brushSize;
    this.ctx.lineJoin = 'round';
    this.ctx.beginPath();
    this.ctx.moveTo(last.x, last.y);
    this.ctx.lineTo(current.x, current.y);
    this.ctx.closePath();
    this.ctx.stroke();
  };
}

export default PaintCanvas;
