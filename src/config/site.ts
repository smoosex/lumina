export interface Project {
  name: string;
  icon: Component;
  link: string;
  color: string;
  external: boolean;
}

export interface Skill {
  name: string;
  icon: Component;
}

export const siteConfig = {
  metadata: {
    title: "Lumina",
    description: "My personal portfolio website",
  },
  profile: {
    name: "Smoose",
    role: "Full Stack Developer",
    email: "kongzhh1101@gmail.com",
    github: "https://github.com/smoosex",
    passion: "Building cool stuff",
  },
};

export const aboutCode = (profile: typeof siteConfig.profile) => `package main

import "${"fmt"}"

type Developer struct {
    Name    string
    Role    string
    Skills  []string
    Passion string
}

func main() {
    me := Developer{
        Name:    "${profile.name}",
        Role:    "${profile.role}",
        Skills:  []string{"Vue", "Go", "TS"},
        Passion: "${profile.passion}",
    }

    fmt.Printf("Hello, I'm %s!\\n", me.Name)
    fmt.Println("Ready to code? 🚀")
}`;

// Icons are auto-imported by unplugin-icons
// If you want to add more icons, check https://icones.js.org/
export const skills: Skill[] = [
  { name: "Go", icon: IconSkillIconsGolang },
  { name: "Kafka", icon: IconSkillIconsKafka },
  { name: "MySQL", icon: IconSkillIconsMysqlDark },
  { name: "Mongo", icon: IconSkillIconsMongodb },
  { name: "Redis", icon: IconSkillIconsRedisDark },
  { name: "Vue", icon: IconSkillIconsVuejsDark },
  { name: "HTML", icon: IconSkillIconsHtml },
  { name: "CSS", icon: IconSkillIconsCss },
  { name: "JavaScript", icon: IconSkillIconsJavascript },
  { name: "TypeScript", icon: IconSkillIconsTypescript },
  { name: "Tailwind CSS", icon: IconSkillIconsTailwindcssDark },
  { name: "Python", icon: IconSkillIconsPythonDark },
  { name: "Lua", icon: IconSkillIconsLuaDark },
];

export const works: Project[] = [
  {
    name: "projectAlpha",
    icon: IconLucideLayout,
    link: "#",
    color: "from-chart-1 to-chart-2",
    external: false,
  },
  {
    name: "projectBeta",
    icon: IconLucideZap,
    link: "#",
    color: "from-chart-2 to-chart-3",
    external: false,
  },
  {
    name: "projectGamma",
    icon: IconLucideClover,
    link: "#",
    color: "from-chart-3 to-chart-4",
    external: false,
  },
];
