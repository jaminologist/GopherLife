$(document).ready(function () {

    var ci = new CanvasInformation()
    OpenWorld(ci)

    var $canvas = $("#worldCanvas");

    $canvas.on('wheel', function (event) {
        ScrollCanvas(event, ci)
        event.preventDefault()
    });

    $canvas.on('keydown', function(event) {
        var key = event.which;
        $.ajax({
            type: 'GET',
            url: '/ShiftWorldView?keydown=' + key,
        })
    });

    $canvas.on('click', function(event) {
        HandleClick(event, ci)
    });
})


function ScrollCanvas(event, CanvasInformation){
    if (event.originalEvent.deltaY < 0) { //Wheel Up
        CanvasInformation.TileWidth = CanvasInformation.TileWidth + 1
        CanvasInformation.TileHeight = CanvasInformation.TileHeight + 1
    } else { //Wheel Down
        CanvasInformation.TileWidth -= 1
        CanvasInformation.TileHeight -= 1

        if (CanvasInformation.TileWidth <= 1) {
            CanvasInformation.TileWidth = 1
            CanvasInformation.TileHeight = 1
        }
    }
}


function UpdateWorldDisplay(data, CanvasInformation) {
    $("#worldDiv").html(data.WorldRender)
    DrawGrid(data.Grid, CanvasInformation)
    DisplaySelectedGopher(data.SelectedGopher)
}


function CanvasInformation(){
    this.TileWidth = 5;
    this.TileHeight = 5;
    this.StartX = 0;
    this.StartY = 0;
    this.RenderWidth = 0;
    this.RenderHeight = 0;
    this.Grid = {}
    this.OtherStartX = 0
    this.OtherStartY = 0
}



function OpenWorld(CanvasInformation) {
    $.ajax({
        url: '/ProcessWorld',
        dataType: 'json',
        type: 'GET',
        success: function (data) {
            CanvasInformation.OtherStartX = data.StartX;
            CanvasInformation.OtherStartY = data.StartY;
            UpdateWorldDisplay(data, CanvasInformation);
            OpenWorld(CanvasInformation);
        },
    })
}

function HandleClick(event, CanvasInformation) {
    var canvas = document.querySelector('canvas')
    var rect = canvas.getBoundingClientRect();

    var canvasX = event.clientX - rect.left;
    var canvasY = event.clientY - rect.top;

    var x = Math.ceil((canvasX - CanvasInformation.StartX) / CanvasInformation.TileWidth);
    var y = Math.ceil((canvasY - CanvasInformation.StartY) / CanvasInformation.TileHeight);

    x = (CanvasInformation.OtherStartX + x) - 1
    y = (CanvasInformation.OtherStartY + y) - 1

    $.ajax({
        type: 'GET',
        url: `/Click?x=${x}&y=${y}`,
        contentType: "application/json",
    });

}

function DrawGrid(Grid, CanvasInformation) {

    var canvas = document.querySelector('canvas')
    resizeCanvasToDisplaySize(canvas)
    var c = canvas.getContext('2d');

    c.clearRect(0, 0, canvas.width, canvas.height);

    CanvasInformation.Grid = Grid
    CanvasInformation.RenderWidth = CanvasInformation.TileWidth * Grid.length
    CanvasInformation.RenderHeight = CanvasInformation.TileHeight * Grid[0].length

    CanvasInformation.StartX = (canvas.width - CanvasInformation.RenderWidth) / 2
    CanvasInformation.StartY = (canvas.height - CanvasInformation.RenderHeight) / 2

    for (var i = 0; i < Grid.length; i++) {
        for (var j = 0; j < Grid[i].length; j++) {

            c.fillStyle = `rgba(${Grid[i][j].R}, ${Grid[i][j].G}, ${Grid[i][j].B}, ${Grid[i][j].A})`; 

            var x = CanvasInformation.StartX + (i * CanvasInformation.TileWidth)
            var y = CanvasInformation.StartY + (j * CanvasInformation.TileHeight)
            c.fillRect(x, y, CanvasInformation.TileWidth, CanvasInformation.TileHeight);
        }
    }
}

function resizeCanvasToDisplaySize(canvas) {
    const width = canvas.clientWidth;
    const height = canvas.clientHeight;

    if (canvas.width !== width || canvas.height !== height) {
        canvas.width = width;
        canvas.height = height;
        return true;
    }

    return false;
}


function DisplaySelectedGopher(gopher) {
    $("#gopher-name").html(gopher.Name)

    var x = gopher.Position.X
    var y = gopher.Position.Y

    $("#gopher-position").html("(" + x + "," + y + ")")
    $("#gopher-hunger").html("(" + gopher.Hunger + ")")
    $("#gopher-lifespan").html("(" + gopher.Lifespan + ")")
}