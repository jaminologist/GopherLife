$(document).ready(function(){

    getWorld()

    $(document).keydown(function(e) {
        var key = e.which;

        var leftArrow = 37
        var rightArrow = 39
        var upArrow = 38
        var downArrow = 40

        var inputData = new FormData()
        inputData.append("keydown", key)



        $.ajax({
            type: 'GET',
            url: '/ShiftWorldView?keydown=' + key,
            error: function(xhr, statusText, err) {
            },
            success: function(){
            }
         })
    });


    $(document).on("click", "a.interactable", function(){

        var position = $(this).attr("id")

        $.ajax({
            type: 'GET',
            url: '/SelectGopher?position=' + position,
            dataType: 'json',
            error: function(xhr, statusText, err) {
            },
            success: function(data){
              UpdateWorldDisplay(data)
            }
         })

    });


})

function UpdateWorldDisplay(data){
    $("#worldDiv").html(data.WorldRender)
    $("#gopher-name").html(data.SelectedGopher.Name)

    var x = data.SelectedGopher.Position.X
    var y = data.SelectedGopher.Position.Y

    $("#gopher-position").html("(" + x + "," + y + ")")
}


function getWorld(){

    $.ajax({
        url: '/PollWorld',
        dataType: 'json',
        success: function(data) {
          UpdateWorldDisplay(data)
        },
        type: 'GET'
     }).always(function(){
        setTimeout(getWorld(),5000);
     });

}
