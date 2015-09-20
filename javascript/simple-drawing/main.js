$(function() {
  var $addNewColor = $('#revealColorSelect');
  var $colorSelect = $('#colorSelect');
  var $newColor = $('#newColor');
  var color = $('.selected').css('background-color');
  var $canvas = $("canvas");
  var context = $canvas[0].getContext("2d");
  var lastEvent;
  var mouseDown = false;

  // 显示/隐藏 添加颜色面板
  $addNewColor.on('click', function() {
    $colorSelect.toggle();
  });

  // 切换选中颜色, 支持动态添加的 li
  $('ul').on('click', 'li', function() {
    // 取消原有的选中状态
    // 同级的其他元素
    $(this).siblings().removeClass('selected');
    // 选中当前元素
    $(this).addClass('selected');
    color = $(this).css('background-color');
  });

  // 新增颜色面板
  // range 值更改时，自动更新 newColor 的颜色
  $('input[type=range]').on('change', function() {
    var r = $('#red').val();
    var g = $('#green').val();
    var b = $('#blue').val();
    $newColor.css('background-color', 'rgb(' + r + ',' + g + ',' + b + ')');
  });

  // 将新颜色添加到 list
  $('#addNewColor').on('click', function() {
    var $li = $('<li></li>');
    $li.css('background-color', $newColor.css('background-color'));
    $('.controls ul').append($li);
    // 选中新增的颜色
    $li.click();
  });

  // 画图
  $canvas.mousedown(function(e){  // 鼠标按下
    lastEvent = e;
    mouseDown = true;
  }).mousemove(function(e){       // 鼠标移动
    if(mouseDown) {               // 只鼠标按住的移动有效
      context.beginPath();        // 开始确定一条线
      context.moveTo(lastEvent.offsetX, lastEvent.offsetY);   // 从最后一次的坐标轴开始
      context.lineTo(e.offsetX, e.offsetY);                   // 到达目前鼠标所在位置
      context.strokeStyle = color;                            // 线的颜色
      context.stroke();                                       // 把线画到面板上
      lastEvent = e;
    }
  }).mouseup(function() {         // 鼠标释放
    mouseDown = false;
  }).mouseleave(function() {      // 鼠标离开 canvas 面板
    $canvas.mouseup();
  });

});
