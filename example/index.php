<?php

echo "<link rel='stylesheet' href='/css/style.css' />";

echo "PHOOOOOOOO ;)";

echo sprintf(<<<"HTML"
     <pre style="font-size: 10px !important; text-align: left">
         <code>
             %s
         </code>
     </pre>
HTML, print_r($_SERVER, true));
