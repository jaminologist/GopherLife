var tileWidth = 5
var tileHeight = 5

$(document).ready(function () {


    OpenWorld()


    var $canvas = $("#worldCanvas");

    $canvas.on('wheel', function (event) {

        if (event.originalEvent.deltaY < 0) { //Wheel Up
            tileWidth += 1
            tileHeight += 1
        } else { //Wheel Down
            tileWidth -= 1
            tileHeight -= 1

            if (tileWidth <= 1) {
                tileWidth = 1
                tileHeight = 1
            }
        }
        e.preventDefault()
    });

   /* $(document).on("submit", "form#reset", function (e) {
        var form = $(this);
        var url = form.attr('action');

        $.ajax({
            type: "GET",
            url: url,
            data: form.serialize(),
        });
        e.preventDefault(); 
    }) */

    $(document).keydown(function (e) {
        var key = e.which;
        $.ajax({
            type: 'GET',
            url: '/ShiftWorldView?keydown=' + key,
        })
    });
})

function UpdateWorldDisplay(data) {
    $("#worldDiv").html(data.WorldRender)
    DrawGrid(data.Grid)
    DisplaySelectedGopher(data.SelectedGopher)
}

function OpenWorld() {
    $.ajax({
        url: '/ProcessWorld',
        dataType: 'json',
        type: 'GET',
        success: function (data) {
            UpdateWorldDisplay(data)
            OpenWorld()
        },
    }).always(function () {
        //setTimeout(OpenWorld(), 5000);
    });
}

Grid = function(canvas){

    this.canvas = canvas
    this.draw = function(){

    }


}



function DrawGrid(Grid) {

    var canvas = document.querySelector('canvas')
    resizeCanvasToDisplaySize(canvas)
    var c = canvas.getContext('2d');

    c.clearRect(0, 0, canvas.width, canvas.height);

    var renderWidth = tileWidth * Grid.length
    var renderHeight = tileHeight * Grid[0].length

    var startX = (canvas.width - renderWidth) / 2
    var startY = (canvas.height - renderHeight) / 2

    for (var i = 0; i < Grid.length; i++) {
        for (var j = 0; j < Grid[i].length; j++) {

            //console.log(Grid[i][j])

            c.fillStyle = `rgba(${Grid[i][j].R}, ${Grid[i][j].G}, ${Grid[i][j].B}, ${Grid[i][j].A})`; 

            var x = startX + (i * tileWidth)
            var y = startY + (j * tileHeight)
            c.fillRect(x, y, tileWidth, tileHeight);
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