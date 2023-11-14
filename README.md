# Jetbrains Blog Template

I read about this in a [JetBrains article](https://blog.jetbrains.com/go/2022/11/08/build-a-blog-with-go-templates/#rendering-templates),
but they also have the source code available [here](https://github.com/JetBrains/go-code-samples/tree/main/GoBlog).

The main Go-takeaway is the Chi framework / router. Looks really cool.

Also this is a learning opportunity for me as far as Go templates _go_ (pun intended). Composing multiple of them into a single page is really powerful.
And since this template is by JetBrains, I guess they know a thing or two about how these kinds of projects go.

## Tailwind

Running tailwind in `watch` mode keeps it refreshing constantly:

```shell
tailwindcss -i ./static/styles.css -o ./static/processed-styles.css --watch
```

The `./static/styles.css` contains hand-written styles that I want to always include with the website.

For production, you'd probably want to run the tailwindcss command with the `--minify` flag too! I just liked seeing the classes being added :)

Adding typography was surprisingly easy - just run `npm i -D @tailwind/typography` as with any old javascript project. I guess the tailwind cli tool will find the packages and resolve them correctly.
