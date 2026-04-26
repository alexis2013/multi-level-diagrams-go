function snapToGrid(value, gridSize) {
    return Math.round(value / gridSize) * gridSize;
}

function drawGrid(layer, stage, gridSize) {
    layer.destroyChildren();
    const width = stage.width();
    const height = stage.height();

    // Vertical lines
    for (let i = 0; i <= width / gridSize; i++) {
        layer.add(new Konva.Line({
            points: [i * gridSize, 0, i * gridSize, height],
            stroke: '#ddd',
            strokeWidth: 0.5,
        }));
    }

    // Horizontal lines
    for (let j = 0; j <= height / gridSize; j++) {
        layer.add(new Konva.Line({
            points: [0, j * gridSize, width, j * gridSize],
            stroke: '#ddd',
            strokeWidth: 0.5,
        }));
    }
}

window.addEventListener('DOMContentLoaded', () => {
    const container = document.getElementById('canvas-container');
    let activeTool = null;
    let gridVisible = false;

    // Initialize Konva Stage
    const stage = new Konva.Stage({
        container: 'canvas-container',
        width: container.offsetWidth,
        height: container.offsetHeight,
    });

    // Add layers (Grid first so it's behind)
    const gridLayer = new Konva.Layer({ visible: false });
    const shapeLayer = new Konva.Layer();
    stage.add(gridLayer);
    stage.add(shapeLayer);

    const getDragBoundFunc = () => {
        return function(pos) {
            if (!gridVisible) return pos;
            return {
                x: snapToGrid(pos.x, GRID_SIZE),
                y: snapToGrid(pos.y, GRID_SIZE)
            };
        };
    };

    // Shape factory
    const createShape = (tool) => {
        let shape;
        const defaults = { ...DEFAULTS[tool] };
        if (!defaults) return;

        if (gridVisible) {
            defaults.x = snapToGrid(defaults.x, GRID_SIZE);
            defaults.y = snapToGrid(defaults.y, GRID_SIZE);
        }

        switch (tool) {
            case 'rect':
                shape = new Konva.Rect({
                    ...defaults,
                    draggable: true,
                    dragBoundFunc: getDragBoundFunc()
                });
                break;
            case 'ellipse':
                shape = new Konva.Ellipse({
                    x: defaults.x,
                    y: defaults.y,
                    radiusX: defaults.width / 2,
                    radiusY: defaults.height / 2,
                    fill: defaults.fill,
                    stroke: defaults.stroke,
                    strokeWidth: defaults.strokeWidth,
                    draggable: true,
                    dragBoundFunc: getDragBoundFunc()
                });
                break;
            case 'text':
                shape = new Konva.Text({
                    ...defaults,
                    draggable: true,
                    dragBoundFunc: getDragBoundFunc()
                });
                break;
            case 'arrow':
                shape = new Konva.Arrow({
                    x: defaults.x,
                    y: defaults.y,
                    points: [0, 0, defaults.width, 0],
                    pointerLength: defaults.pointerLength,
                    pointerWidth: defaults.pointerWidth,
                    fill: defaults.fill,
                    stroke: defaults.stroke,
                    strokeWidth: defaults.strokeWidth,
                    draggable: true,
                    dragBoundFunc: getDragBoundFunc()
                });
                break;
        }

        if (shape) {
            shapeLayer.add(shape);
            shapeLayer.draw();
        }
    };

    // Tool selection
    const buttons = document.querySelectorAll('.sidebar button[data-tool]');
    buttons.forEach(btn => {
        btn.addEventListener('click', () => {
            const tool = btn.getAttribute('data-tool');
            buttons.forEach(b => b.classList.remove('active'));
            
            if (tool === 'select') {
                activeTool = null;
            } else {
                activeTool = tool;
                btn.classList.add('active');
                createShape(tool);
            }
        });
    });

    // Grid toggle
    const gridToggle = document.getElementById('grid-toggle');
    gridToggle.addEventListener('click', () => {
        gridVisible = !gridVisible;
        gridLayer.visible(gridVisible);
        if (gridVisible) {
            drawGrid(gridLayer, stage, GRID_SIZE);
            gridToggle.textContent = 'Hide Grid';
        } else {
            gridToggle.textContent = 'Show Grid';
        }
        gridLayer.draw();
    });

    // Handle window resize
    window.addEventListener('resize', () => {
        const w = container.offsetWidth;
        const h = container.offsetHeight;
        stage.width(w);
        stage.height(h);
        if (gridVisible) {
            drawGrid(gridLayer, stage, GRID_SIZE);
        }
        stage.draw();
    });
});
