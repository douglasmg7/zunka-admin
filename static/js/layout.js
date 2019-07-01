(function (window, document) {

    // Hamburg.
    var menu   = document.getElementById('menu');
    var burger = document.getElementsByClassName('burger')[0];


    // // media query event handler
    // if (matchMedia) {
    //     const mq = window.matchMedia("(min-width: 48em)");
    //     // const mq = window.matchMedia("(min-width: 500px)");
    //     mq.addListener(widthChange);
    //     widthChange(mq);
    // }

    // // media query change
    // function widthChange(mq) {
    //     // More then 48em.
    //     if (mq.matches) {
    //         console.log('1')
    //         let toggleItems = document.getElementsByClassName("toggle-items");
    //         for (let i = 0; i < toggleItems.length; i++) {
    //             let subMenu = toggleItems[i].nextElementSibling;
    //             console.log(toggleItems[i]);
    //             console.log(toggleItems[i].offsetLeft);
    //             console.log(subMenu);

    //             subMenu.style.left = toggleItems[i].offsetLeft  + "px";
    //         }
    //     // Less than 48em.
    //     } else {
    //         console.log('2')
    //     }
    // }


    // Toggle all.
    function toggleMenu(e) {
        e.preventDefault();
        toggleClass(menu, "active");
        toggleClass(burger, "active");
    }

    // Show and hide menu.
    burger.onclick = function (e) {
        toggleMenu(e);
    };

    // Show and hide sub-menus. 
    let toggleItems = document.getElementsByClassName("toggle-items");
    for (let i = 0; i < toggleItems.length; i++) {
        let subMenu = toggleItems[i].nextElementSibling;
        // Show sub-menu.
        toggleItems[i].onclick =  function(e){
            toggleClass(subMenu, "active");
            // subMenu.style.top = "200px";
            // subMenu.style.left = "200px";
        }
        // Hide sub-menu.
        let backElements = subMenu.getElementsByClassName("back");
        backElements[0].onclick = function(e){
            toggleClass(subMenu, "active");
        }
    }


    // Toggle class.
    function toggleClass(element, className) {
        var classes = element.className.split(/\s+/),
            length = classes.length,
            i = 0;

        for(; i < length; i++) {
          if (classes[i] === className) {
            classes.splice(i, 1);
            break;
          }
        }
        // The className is not found
        if (length === classes.length) {
            classes.push(className);
        }

        element.className = classes.join(' ');
    }


}(this, this.document));
