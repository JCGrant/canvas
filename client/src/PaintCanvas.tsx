import * as React from 'react';

const canvasStyle = {
  border: '1px solid black',
};

interface IPaintCanvasState {
  brushColor: string;
  brushSize: number;
  mouseDown: boolean;
}

class PaintCanvas extends React.Component<{}, IPaintCanvasState> {
  private canvas: HTMLCanvasElement;
  private ctx: CanvasRenderingContext2D;

  constructor(props: {}) {
    super(props);
    this.state = {
      brushColor: '#000000',
      brushSize: 10,
      mouseDown: false,
    };
  }

  public render() {
    return (
      <div>
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

  public componentWillUnmount() {
    window.removeEventListener('mouseup', this.onMouseUp);
  }

  private canvasRef = (canvas: HTMLCanvasElement) => {
    this.canvas = canvas;
    this.ctx = canvas.getContext('2d')!;
  };

  private onMouseMove = (e: React.MouseEvent) => {
    if (!this.state.mouseDown) {
      return;
    }
    this.paint(
      e.pageX - this.canvas.offsetLeft,
      e.pageY - this.canvas.offsetTop,
      this.state.brushSize,
      this.state.brushColor
    );
  };

  private onMouseDown = (e: React.MouseEvent) => {
    e.persist();
    this.setState({ mouseDown: true }, () => this.onMouseMove(e));
  };

  private onMouseUp = () => {
    this.setState({ mouseDown: false });
  };

  private paint = (
    x: number,
    y: number,
    brushSize: number = 1,
    color: string = 'black'
  ) => {
    this.ctx.fillStyle = color;
    this.ctx.beginPath();
    this.ctx.arc(x, y, brushSize, 0, Math.PI * 2, true);
    this.ctx.fill();
  };
}

export default PaintCanvas;
