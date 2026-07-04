import { Check, Moon, Sun, Trees } from "lucide-react"

import { Button } from "@/components/ui/button"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { useTheme, type Theme } from "@/components/theme-provider"

export function ThemeSwitcher() {
  const { theme, setTheme } = useTheme()

  const currentIcon = () => {
    switch (theme) {
      case "dark":
        return <Moon className="h-[1.2rem] w-[1.2rem] transition-all" />
      case "british":
        return <Trees className="h-[1.2rem] w-[1.2rem] transition-all" />
      default:
        return <Sun className="h-[1.2rem] w-[1.2rem] transition-all" />
    }
  }

  const isCurrent = (themeName: Theme) => {
    // If theme is system, we might not have a direct match without checking window preferences,
    // but typically we can just show what they manually selected.
    return theme === themeName
  }

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="ghost" size="icon" className="focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2">
          {currentIcon()}
          <span className="sr-only">Toggle theme</span>
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end">
        <DropdownMenuItem onClick={() => setTheme("light")} className="flex items-center justify-between">
          Light
          {isCurrent("light") && <Check className="h-4 w-4" />}
        </DropdownMenuItem>
        <DropdownMenuItem onClick={() => setTheme("dark")} className="flex items-center justify-between">
          Dark
          {isCurrent("dark") && <Check className="h-4 w-4" />}
        </DropdownMenuItem>
        <DropdownMenuItem onClick={() => setTheme("british")} className="flex items-center justify-between">
          British
          {isCurrent("british") && <Check className="h-4 w-4" />}
        </DropdownMenuItem>
        <DropdownMenuItem onClick={() => setTheme("system")} className="flex items-center justify-between">
          System
          {isCurrent("system") && <Check className="h-4 w-4" />}
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  )
}
